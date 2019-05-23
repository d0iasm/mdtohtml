package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

type Type int

const (
	RAWTEXT Type = iota
	BODY
	HEADING
	H1
	H2
	H3
	P
	LIST
	LINK
	BR
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fname := os.Args[1]
	name := strings.Split(fname, ".")

	if strings.Compare(strings.ToLower(name[len(name)-1]), "md") != 0 {
		fmt.Println("Input file must be a markdown file(.md).")
		os.Exit(1)
	}

	wfile, err := os.Create(name[0] + ".html")
	check(err)
	defer wfile.Close()
	writer := bufio.NewWriter(wfile)
	if len(os.Args) < 3 || os.Args[2] != "-nocss" {
		fmt.Fprintln(writer, css())
	}

	rfile, err := os.Open(fname)
	check(err)
	defer rfile.Close()
	reader := bufio.NewReader(rfile)

	t := Tokenizer{bufio.NewScanner(reader), bytes.NewBufferString("")}
	tokens := t.tokenize()
	fmt.Println("TOKENS: ", tokens)

	p := Parser{0, tokens}
	root := p.body()
	debugTree(root, 0)
	fmt.Println("NODES: ", root)

	html := generate(root)
	//fmt.Println("HTML: ", html)

	fmt.Fprintln(writer, html)
	writer.Flush()
}
