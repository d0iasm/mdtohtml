package main

import (
	"bufio"
	//"fmt"
	"strconv"
	"strings"
)

var headOfLine = true

type Type int

const (
	RAWTEXT Type = iota
	BODY
	P
	HEADING
	UL
	LIST
	LINK
	URI
	EOF
)

var block = []Type{P, HEADING, UL}
var inline = []Type{RAWTEXT, LINK}

type Token struct {
	ty  Type
	val string
	dep int
}

type Tokenizer struct {
	s   *bufio.Scanner
	buf string
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

func (t *Tokenizer) heading() {
	n := t.count("#")
	if n > 6 {
		t.buf += strings.Repeat("#", n) + t.s.Text()
		return
	}

	if t.s.Text() == " " {
		t.tokens = append(t.tokens, Token{HEADING, strings.Repeat("#", n), -1})
		t.inline()
		return
	}
	t.buf += strings.Repeat("#", n) + t.s.Text()
}

func (t *Tokenizer) ul(dep int, sym string) {
	t.tokens = append(t.tokens, Token{UL, "", dep})
	t.list(dep, sym)
	// Create an unnested list with this depth.
	if t.buf == strings.Repeat(" ", (dep*2))+sym && t.s.Text() == " " {
		headOfLine = true
		t.buf = ""
		t.list(dep, sym)
	}
}

func (t *Tokenizer) list(dep int, sym string) {
	t.tokens = append(t.tokens, Token{LIST, sym, dep})

	t.listitem()
	if !t.s.Scan() {
		return
	}

	// |headOfLine| should be false when heading() is called.
	headOfLine = true

	// Check an existance of an unnested list.
	for i := 0; i < dep*2; i++ {
		t.buf += t.s.Text()
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
			t.buf = ""
			t.list(dep, sym)
		} else {
			t.buf += sym
		}
	case " ": // Create a nested sublist.
		n := t.count(" ")
		t.buf += strings.Repeat(" ", n)
		if n == 2 {
			t.buf = ""
			if t.s.Text() != sym {
				t.buf += " " + t.s.Text()
				return
			}
			if !t.consume(" ") {
				t.buf += " " + sym + t.s.Text()
				return
			}
			t.ul(dep+1, sym)
		}
	}
}

func (t *Tokenizer) listitem() {
	for t.s.Scan() {
		switch t.s.Text() {
		case "\n":
			t.tokens = append(t.tokens, Token{RAWTEXT, t.buf, -1})
			t.buf = ""
			return
		case "#":
			if len(t.buf) == 0 {
				t.heading()
			} else {
				t.buf += t.s.Text()
			}
		case "[":
			if len(t.buf) > 0 {
				t.tokens = append(t.tokens, Token{RAWTEXT, t.buf, -1})
				t.buf = ""
			}
			t.link()
		default:
			t.buf += t.s.Text()
		}
	}
	if len(t.buf) > 0 {
		t.tokens = append(t.tokens, Token{RAWTEXT, t.buf, -1})
		t.buf = ""
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
				tmp := t.buf
				t.buf = ""
				t.buf += "[" + tmp
				return
			}
			if posUriS < 0 {
				tmp := t.buf
				t.buf = ""
				t.buf += "[" + link + "]" + tmp
				return
			}
			if posUriE < 0 {
				tmp := t.buf
				t.buf = ""
				t.buf += "[" + link + "]" + "(" + tmp
			}
			return
		case "]":
			link = t.buf
			t.buf = ""
			posLinkE = i
		case "(":
			if posLinkE < 0 {
				t.buf += "("
				break
			}

			if posLinkE+1 != i {
				tmp := t.buf
				t.buf = ""
				t.buf += "[" + link + "]" + tmp + "("
				return
			}
			posUriS = i
		case ")":
			if posLinkE < 0 || posUriS < 0 {
				t.buf += ")"
				break
			}

			posUriE = i
			if posLinkE+1 == posUriS && posUriS < posUriE {
				uri = t.buf
				t.buf = ""
				t.tokens = append(t.tokens, Token{LINK, link, -1})
				t.tokens = append(t.tokens, Token{URI, uri, -1})
				return
			}
		default:
			t.buf += t.s.Text()
		}
		i++
	}
}

func (t *Tokenizer) inline() {
	for t.s.Scan() {
		switch t.s.Text() {
		case "[":
			if len(t.buf) > 0 {
				t.tokens = append(t.tokens, Token{RAWTEXT, t.buf, -1})
				t.buf = ""
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
		switch t.s.Text() {
		case "\n":
			if len(t.buf) <= 0 {
				t.buf += " "
				break
			}
			if headOfLine {
				t.tokens = append(t.tokens, Token{P, t.buf, -1})
			} else {
				t.tokens = append(t.tokens, Token{RAWTEXT, t.buf, -1})
			}
			t.buf = ""
			headOfLine = true
		case "#":
			if !headOfLine {
				t.buf += t.s.Text()
				break
			}
			t.heading()
		case "-":
			sym := t.s.Text()
			if t.consume(" ") {
				t.ul(0, sym)
				break
			}
			t.buf += sym + t.s.Text()
		case "[":
			if len(t.buf) > 0 {
				t.tokens = append(t.tokens, Token{RAWTEXT, t.buf, -1})
				t.buf = ""
			}
			t.link()
		default:
			t.buf += t.s.Text()
		}
	}
	if len(t.buf) > 0 {
		t.tokens = append(t.tokens, Token{P, t.buf, -1})
	}
	t.tokens = append(t.tokens, Token{EOF, "", -1})
}

func (t *Tokenizer) getTokens() []Token {
	return t.tokens
}
