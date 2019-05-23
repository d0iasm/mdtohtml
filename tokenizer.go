package main

import (
	"bufio"
	"bytes"
	//"fmt"
	"strings"
)

type Token struct {
	ty  Type
	val string
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
			return strings.Join(buf, "")
		}
		buf = append(buf, t.s.Text())
	}
	return strings.Join(buf, "")
}

func (t *Tokenizer) consume(target string) bool {
	if t.s.Scan() && t.s.Text() == target {
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

func (t *Tokenizer) ul(nest int, sym string) {
	t.tokens = append(t.tokens, Token{UL, ""})
	t.list(nest, sym)
	t.tokens = append(t.tokens, Token{UL_END, ""})
	if t.buf.String() == sym+" " {
		t.buf.Reset()
		t.list(nest, sym)
	}
}

func (t *Tokenizer) list(nest int, sym string) {
	t.tokens = append(t.tokens, Token{LIST, sym})
	t.tokens = append(t.tokens, Token{RAWTEXT, t.stringLiteral()})
	t.consume("\n")

	for i := 0; i < (nest * 2); i++ {
		t.buf.WriteString(t.s.Text())
		if !t.s.Scan() {
			return
		}
	}

	switch t.s.Text() {
	case sym:
		if t.consume(" ") {
			t.list(nest, sym)
		}
	case " ":
		n := t.count(" ")
		t.buf.WriteString(strings.Repeat(" ", n))
		if n == (nest+1)*2 {
			t.buf.Reset()
			t.ul(nest+1, sym)
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
				t.tokens = append(t.tokens, Token{P, t.buf.String()})
			} else {
				t.tokens = append(t.tokens, Token{RAWTEXT, t.buf.String()})
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
				t.tokens = append(t.tokens, Token{HEADING, strings.Repeat("#", n)})
				// TODO: Not only rawtext after heading. such as link...
				isHead = false
				break
			}
			t.buf.WriteString(strings.Repeat("#", n))
			t.buf.WriteString(t.s.Text())
		case "-":
			sym := t.s.Text()
			if t.consume(" ") {
				t.ul(0, sym)
				break
			}
			t.buf.WriteString(sym)
			t.buf.WriteString(t.s.Text())
		default:
			t.buf.WriteString(t.s.Text())
		}
	}
}

func (t *Tokenizer) getTokens() []Token {
	return t.tokens
}
