package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

var headOfLine = true

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

func (t *Tokenizer) heading() {
	n := t.count("#")
	if n > 6 {
		t.buf.WriteString(strings.Repeat("#", n))
		t.buf.WriteString(t.s.Text())
		return
	}

	if t.s.Text() == " " {
		t.tokens = append(t.tokens, Token{HEADING, strings.Repeat("#", n), -1})
		t.inline()
		return
	}
	t.buf.WriteString(strings.Repeat("#", n))
	t.buf.WriteString(t.s.Text())
}

func (t *Tokenizer) ul(dep int, sym string) {
	t.tokens = append(t.tokens, Token{UL, "", dep})
	t.list(dep, sym)
	// Create an unnested list with this depth.
	if t.buf.String() == strings.Repeat(" ", (dep*2))+sym && t.s.Text() == " " {
		headOfLine = true
		t.buf.Reset()
		t.list(dep, sym)
	}
}

func (t *Tokenizer) list(dep int, sym string) {
	t.tokens = append(t.tokens, Token{LIST, sym, dep})
	t.inline()

	// End of a list. Return when a fist char is block element.
	if t.s.Text() == "\n" {
		return
	}
	if t.s.Text() == "#" {
		t.heading()
		return
	}

	// |headOfLine| should be false when heading() is called.
	headOfLine = true

	// Check an existance of an unnested list.
	for i := 0; i < dep*2; i++ {
		t.buf.WriteString(t.s.Text())
		if t.s.Text() == "-" && t.consume(" ") {
			return
		}
		if !t.s.Scan() {
			return
		}
	}

	switch t.s.Text() {
	case sym: // Continue a list with the same depth.
		if t.consume(" ") {
			t.buf.Reset()
			t.list(dep, sym)
		} else {
			t.buf.WriteString(sym)
		}
	case " ": // Create a nested sublist.
		n := t.count(" ")
		t.buf.WriteString(strings.Repeat(" ", n))
		if n == 2 {
			t.buf.Reset()
			if t.s.Text() != sym {
				t.buf.WriteString(" ")
				t.buf.WriteString(t.s.Text())
				return
			}
			if !t.consume(" ") {
				t.buf.WriteString(" ")
				t.buf.WriteString(sym)
				t.buf.WriteString(t.s.Text())
				return
			}
			t.ul(dep+1, sym)
		}
	}
}

func (t *Tokenizer) link() {
	posLinkE := -1
	posUriS := -1
	posUriE := -1
	i := 0
	link := ""
	uri := ""
	for t.s.Scan() {
		switch t.s.Text() {
		case "\n":
			if posLinkE < 0 {
				tmp := t.buf.String()
				t.buf.Reset()
				t.buf.WriteString("[" + tmp)
				return
			}
			if posUriS < 0 {
				tmp := t.buf.String()
				t.buf.Reset()
				t.buf.WriteString("[" + link + "]" + tmp)
				return
			}
			if posUriE < 0 {
				tmp := t.buf.String()
				t.buf.Reset()
				t.buf.WriteString("[" + link + "]" + "(" + tmp)
			}
			return
		case "]":
			link = t.buf.String()
			t.buf.Reset()
			posLinkE = i
		case "(":
			if posLinkE < 0 {
				t.buf.WriteString("(")
				break
			}

			if posLinkE+1 != i {
				tmp := t.buf.String()
				t.buf.Reset()
				t.buf.WriteString("[" + link + "]" + tmp + "(")
				return
			}
			posUriS = i
		case ")":
			if posLinkE < 0 || posUriS < 0 {
				t.buf.WriteString(")")
				break
			}

			posUriE = i
			if posLinkE+1 == posUriS && posUriS < posUriE {
				uri = t.buf.String()
				t.buf.Reset()
				t.tokens = append(t.tokens, Token{LINK, link, -1})
				t.tokens = append(t.tokens, Token{URI, uri, -1})
				return
			}
		default:
			t.buf.WriteString(t.s.Text())
		}
		i++
	}
}

func (t *Tokenizer) inline() {
	for t.s.Scan() {
		switch t.s.Text() {
		case "[":
			if t.buf.Len() > 0 {
				t.tokens = append(t.tokens, Token{RAWTEXT, t.buf.String(), -1})
				t.buf.Reset()
			}
			t.link()
		default:
                        // TODO: remove stringLiteral().
			t.tokens = append(t.tokens, Token{RAWTEXT, t.stringLiteral(), -1})
		}
		return
	}
}

func (t *Tokenizer) tokenize() {
	t.s.Split(bufio.ScanRunes)
	for t.s.Scan() {
		fmt.Println("before switch:", t.s.Text(), headOfLine)
		switch t.s.Text() {
		case "\n":
			if t.buf.Len() <= 0 {
				t.buf.WriteString(" ")
				break
			}
			if headOfLine {
				t.tokens = append(t.tokens, Token{P, t.buf.String(), -1})
			} else {
				t.tokens = append(t.tokens, Token{RAWTEXT, t.buf.String(), -1})
			}
			t.buf.Reset()
			headOfLine = true
		case "#":
			if !headOfLine {
				t.buf.WriteString(t.s.Text())
				break
			}
			t.heading()
		case "-":
			sym := t.s.Text()
			if t.consume(" ") {
				t.ul(0, sym)
				break
			}
			t.buf.WriteString(sym)
			t.buf.WriteString(t.s.Text())
		case "[":
			if t.buf.Len() > 0 {
				t.tokens = append(t.tokens, Token{RAWTEXT, t.buf.String(), -1})
				t.buf.Reset()
			}
			t.link()
		default:
			t.buf.WriteString(t.s.Text())
		}
		fmt.Println("after switch:", t.s.Text(), headOfLine)
	}
	t.tokens = append(t.tokens, Token{EOF, "", -1})
}

func (t *Tokenizer) getTokens() []Token {
	return t.tokens
}
