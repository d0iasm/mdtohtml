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
	rawtext Type = iota
	h1
	h2
	h3
	li
	p
	br
)

type Token struct {
	ty  Type
	val string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getRawtext(i *int, chars []rune) string {
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
			tokens = append(tokens, Token{h1, string(chars[i])})
			i++
			tokens = append(tokens, Token{rawtext, getRawtext(&i, chars)})
			fmt.Println("Called #", tokens[len(tokens)-1])
                case "\n":
			fmt.Println("Called br", tokens[len(tokens)-1])
                        i++
		default:
			tokens = append(tokens, Token{p, getRawtext(&i, chars)})
			i++
			fmt.Println("Called default", tokens[len(tokens)-1])
		}
	}
	return tokens
}

func generate(tokens []Token) string {
	html := ""
	for i, t := range tokens {
          switch t.ty {
          case h1:
            html += "<h1>" + tokens[i].val + "</h1>"
          case p:
            html += "<p>" + tokens[i].val + "</p>"
          default:
            html += t.val
          }
	}
	return html
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
	html := generate(tokens)

	fmt.Fprintln(writer, html)
	writer.Flush()
}
