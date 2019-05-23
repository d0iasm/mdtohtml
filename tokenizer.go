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
	if t.buf.String() == strings.Repeat(" ", (dep*2))+sym+" " {
                fmt.Println("FIND:", sym+" ", dep)
		t.buf.Reset()
		  t.list(dep, sym)
	}
}

func (t *Tokenizer) list(dep int, sym string) {
	t.tokens = append(t.tokens, Token{LIST, sym, dep})
	t.tokens = append(t.tokens, Token{RAWTEXT, t.stringLiteral(), dep})

        fmt.Println("==============:")
	for i := 0; i < (dep * 2); i++ {
		t.buf.WriteString(t.s.Text())
                fmt.Println("buf:", t.buf.String(), dep, i)
                fmt.Println("target:", strings.Repeat(" ", i)+sym+" ", dep, i)
	        if t.buf.String() == strings.Repeat(" ", i)+sym+" " {
                fmt.Println("FIND:", sym+" ", dep, i)
		t.buf.Reset()
                  t.consume(" ")
		  t.list(i-1, sym)
	        }
		if !t.s.Scan() {
			return
		}
	}

	switch t.s.Text() {
	case sym:
		if t.consume(" ") {
			t.buf.Reset()
			t.list(dep, sym)
		}
	case " ":
		n := t.count(" ")
		t.buf.WriteString(strings.Repeat(" ", n))
		if n == 2 {
                        if t.s.Text() != sym {
                          t.buf.WriteString(t.s.Text())
                          return
                        }
                        if !t.consume(" ") {
                          t.buf.WriteString(t.s.Text())
                          return
                        }
			t.buf.Reset()
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

        t.tokens = append(t.tokens, Token{EOF, "", -1})
}

func (t *Tokenizer) getTokens() []Token {
	return t.tokens
}
