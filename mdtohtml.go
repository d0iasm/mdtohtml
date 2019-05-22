package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
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

	rfile, err := os.Open(fname)
	check(err)
	defer rfile.Close()

	wfile, err := os.Create(name[0] + ".html")
	check(err)
	defer wfile.Close()

	reader := bufio.NewReader(rfile)
	mdbytes, err := ioutil.ReadAll(reader)
	mdchars := bytes.Runes(mdbytes)
	check(err)

	writer := bufio.NewWriter(wfile)
	if len(os.Args) < 3 || os.Args[2] != "-nocss" {
		fmt.Fprintln(writer, css())
	}

	rfile.Seek(0, io.SeekStart)
	reader = bufio.NewReader(rfile)
	s := bufio.NewScanner(reader)

	t := Tokenizer{0, mdchars, s, []Token{}}
	tokens := t.tokenize()
	fmt.Println("TOKENS: ", tokens)

	p := Parser{0, tokens}
	root := p.body()
	//fmt.Println("NODES: ", root)

	html := generate(root)
	//fmt.Println("HTML: ", html)

	fmt.Fprintln(writer, html)
	writer.Flush()
}
