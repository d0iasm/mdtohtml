package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Type int

const (
	RAWTEXT Type = iota
	BODY
	H1
	H2
	H3
	LI
	P
	BR
)

type Token struct {
	ty  Type
	val string
}

type Node struct {
	ty       Type
	children []Node
	val      string
}

type Parser struct {
	i      int
	tokens []Token
}

func appendChild(parent *Node, child Node) {
	parent.children = append(parent.children, child)
}

func (p *Parser) h1() Node {
	t := p.tokens[p.i]
	n := Node{H1, []Node{}, ""}
	p.i++
	t = p.tokens[p.i]
	if t.ty == RAWTEXT {
		appendChild(&n, p.rawtext(t.val))
	}
	return n
}

func (p *Parser) p() Node {
	t := p.tokens[p.i]
	n := Node{P, []Node{}, ""}
	appendChild(&n, p.rawtext(t.val))
	return n
}

func (p *Parser) rawtext(s string) Node {
	return Node{RAWTEXT, []Node{}, s}
}

func (p *Parser) html() Node {
	root := Node{BODY, []Node{}, ""}
	for p.i < len(p.tokens) {
		t := p.tokens[p.i]
		switch t.ty {
		case H1:
			appendChild(&root, p.h1())
		case P:
			appendChild(&root, p.p())
		}
		p.i++
	}
	return root
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func text(i *int, chars []rune) string {
	start := *i
	for string(chars[*i]) != "\n" {
		*i++
		if *i >= len(chars) {
			break
		}
	}
	return string(chars[start:*i])
}

func tokenize(chars []rune) []Token {
	fmt.Println(string(chars))
	tokens := []Token{}
	i := 0
	for i < len(chars) {
		switch string(chars[i]) {
		case "#":
			tokens = append(tokens, Token{H1, "#"})
			i += 2
			tokens = append(tokens, Token{RAWTEXT, text(&i, chars)})
		case "\n":
			i++
		default:
			tokens = append(tokens, Token{P, text(&i, chars)})
			i++
		}
	}
	return tokens
}

func generate(node Node) string {
	fmt.Println("Current Node", node.ty, node.val)
	if node.ty == RAWTEXT {
		return node.val
	}

	html := ""
	for _, c := range node.children {
		html += generate(c)
	}

	switch node.ty {
	case H1:
		return "<h1>" + html + "</h1>"
	case P:
		return "<p>" + html + "</p>"
	default:
		return html
	}
}

func css() string {
	return `
<link href="https://fonts.googleapis.com/css?family=Abril+Fatface|Lora|Noto+Serif|Noto+Serif+JP" rel="stylesheet">
<style>
body {
  font-family: 'Lora', 'Noto Serif', 'Noto Serif JP', serif;
  max-width: 740px;
  margin: 0 auto;
}
h1, h2, h3 {
  font-family: 'Abril Fatface', cursive;
}
</style>
`
}

func main() {
	fname := os.Args[1]
	name := strings.Split(fname, ".")

	if strings.Compare(strings.ToLower(name[1]), "md") != 0 {
		fmt.Println("Input file must be a markdown file(.md).")
		os.Exit(1)
	}

	rfile, err := os.Open(fname)
	check(err)
	defer rfile.Close()

	wfile, err := os.Create(name[0] + ".html")
	check(err)
	defer wfile.Close()

	writer := bufio.NewWriter(wfile)
	fmt.Fprintln(writer, css())

	reader := bufio.NewReader(rfile)
	mdbytes, err := ioutil.ReadAll(reader)
	mdchars := bytes.Runes(mdbytes)
	check(err)

	tokens := tokenize(mdchars)
	fmt.Println("TOKENS: ", tokens)

	parser := &Parser{0, tokens}
	root := parser.html()
	fmt.Println("NODES: ", root)

	html := generate(root)
	fmt.Println("HTML: ", html)

	fmt.Fprintln(writer, html)
	writer.Flush()
}
