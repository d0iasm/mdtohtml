package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

type Token struct {
	ty  Type
	val string
	dep int
}

type Tokenizer struct {
	s      *bufio.Scanner
	buf    *bytes.Buffer
	tokens []Token
}

func (t *Tokenizer) stringLiteral() string {
	buf := []string{t.s.Text()}
	for t.s.Scan() {
		if t.s.Text() == "\n" {
			t.s.Scan()
			return strings.Join(buf, "")
		}
		buf = append(buf, t.s.Text())
	}
	return strings.Join(buf, "")
}

func (t *Tokenizer) consume(target string) bool {
	if t.s.Scan() && t.s.Text() == target {
		t.s.Scan()
		return true
	}
	return false
}

func (t *Tokenizer) count(target string) int {
	n := 1
	for t.s.Scan() {
		if t.s.Text() == target {
			n++
			continue
		}
		return n
	}
	return n
}

func (t *Tokenizer) ul(dep int, sym string) {
	t.tokens = append(t.tokens, Token{UL, "", dep})
	t.list(dep, sym)
	// Check if an unnested list exists or not.
	if t.buf.String() == strings.Repeat(" ", (dep*2))+sym+" " {
		fmt.Println("ul0", t.buf.String())
		t.buf.Reset()
		t.list(dep, sym)
	}
}

func (t *Tokenizer) list(dep int, sym string) {
	t.tokens = append(t.tokens, Token{LIST, sym, dep})
	t.tokens = append(t.tokens, Token{RAWTEXT, t.stringLiteral(), dep})

	// Move whitespaces to a buffer.
	for i := 0; i < (dep * 2); i++ {
		fmt.Println("li0", i, t.buf.String())
		t.buf.WriteString(t.s.Text())
		if t.buf.String() == strings.Repeat(" ", i*2)+sym {
			t.buf.Reset()
			if t.consume(" ") {
				t.list(i, sym)
			} else {
				t.buf.WriteString(sym)
			}
		}
		if !t.s.Scan() {
			return
		}
	}

	fmt.Println("li1", t.buf.String())
	switch t.s.Text() {
	case sym:
		fmt.Println("li2", t.buf.String())
		if t.consume(" ") {
			t.buf.Reset()
			t.list(dep, sym)
		} else {
			t.buf.WriteString(sym)
		}
	case " ":
		fmt.Println("li3", t.buf.String())
		n := t.count(" ")
		t.buf.WriteString(strings.Repeat(" ", n))
		if n == 2 {
			t.buf.Reset()
			if t.s.Text() != sym {
				fmt.Println("li4", t.buf.String())
				t.buf.WriteString(" ")
				t.buf.WriteString(t.s.Text())
				return
			}
			if !t.consume(" ") {
				t.buf.WriteString(" ")
				t.buf.WriteString(sym)
				t.buf.WriteString(t.s.Text())
				fmt.Println("li5", t.buf.String(), t.s.Text())
				return
			}
			fmt.Println("li6", t.buf.String(), t.s.Text())
			t.ul(dep+1, sym)
		}
	}
}

func (t *Tokenizer) tokenize() {
	isHead := true
	t.s.Split(bufio.ScanRunes)
	for t.s.Scan() {
		switch t.s.Text() {
		case "\n":
			if t.buf.Len() <= 0 {
				break
			}
			if isHead {
				t.tokens = append(t.tokens, Token{P, t.buf.String(), -1})
			} else {
				t.tokens = append(t.tokens, Token{RAWTEXT, t.buf.String(), -1})
			}
			t.buf.Reset()
			isHead = true
		case "#":
			if !isHead {
				t.buf.WriteString(t.s.Text())
				break
			}

			n := t.count("#")
			if n > 6 {
				t.buf.WriteString(strings.Repeat("#", n))
				t.buf.WriteString(t.s.Text())
				break
			}

			if t.s.Text() == " " {
				t.tokens = append(t.tokens, Token{HEADING, strings.Repeat("#", n), -1})
				// TODO: Not only rawtext after heading. such as link...
				isHead = false
				break
			}
			t.buf.WriteString(strings.Repeat("#", n))
			t.buf.WriteString(t.s.Text())
		case "-":
			sym := t.s.Text()
			fmt.Println("main", t.buf.String(), sym)
			if t.consume(" ") {
				fmt.Println("main2", t.buf.String(), sym)
				t.ul(0, sym)
				isHead = false
				break
			}
			t.buf.WriteString(sym)
			t.buf.WriteString(t.s.Text())
		default:
			t.buf.WriteString(t.s.Text())
		}
	}

	t.tokens = append(t.tokens, Token{EOF, "", -1})
}

func (t *Tokenizer) getTokens() []Token {
	return t.tokens
}
