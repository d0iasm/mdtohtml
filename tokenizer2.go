package main

import (
	"bufio"
	//"fmt"
	"strings"
)

type TokenType int

const (
	TK_TEXT = iota
	TK_RESERVED
	TK_NEWLINE
	TK_EOF
)

type Token struct {
	ty  TokenType
	val string
	dep int
}

type Tokenizer struct {
	s      *bufio.Scanner
	buf    string
	tokens []Token
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
		case "\n":
			t.tokens = append(t.tokens, Token{RAWTEXT, t.buf, -1})
			t.buf = ""
			return
		default:
			t.buf += t.s.Text()
		}
	}

	if len(t.buf) > 0 {
		t.tokens = append(t.tokens, Token{RAWTEXT, t.buf, -1})
		t.buf = ""
	}
}

func (t *Tokenizer) tokenize() {
	t.s.Split(bufio.ScanRunes)
	headOfLine = true
	for t.s.Scan() {
		switch t.s.Text() {
		case "\n":
			if headOfLine {
				t.tokens = append(t.tokens, Token{P, t.buf, -1})
			} else {
				t.buf += " "
				t.tokens = append(t.tokens, Token{RAWTEXT, t.buf, -1})
			}
			t.buf = ""
			headOfLine = true
		case "#":
			if headOfLine {
				t.heading()
				break
			}
			t.buf += t.s.Text()
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
			headOfLine = false
		}
	}

	if len(t.buf) > 0 {
		t.tokens = append(t.tokens, Token{P, t.buf, -1})
	}
	t.tokens = append(t.tokens, Token{EOF, "", -1})
}

func (t *Tokenizer) startReserved() string {
	keywords := []string{"return", "if", "else", "for", "func", "var", "package"}
	for _, kw := range keywords {
		if strings.HasPrefix(in, kw) {
			if len(kw) == len(in) || !isAlnum(in[len(kw)]) {
				return kw
			}
		}
	}
	return ""
}

func (t *Tokenizer) tokenize() {
	t.s.Split(bufio.ScanRunes)
	for t.s.Scan() {
		switch t.s.Text() {
		case "\n":
			t.tokens = append(t.tokens, Token{TK_NEWLINE, t.s.Text(), -1})
		}
	}

	if len(t.buf) > 0 {
		t.tokens = append(t.tokens, Token{TK_TEXT, t.buf, -1})
	}
	t.tokens = append(t.tokens, Token{EOF, "", -1})
}

func (t *Tokenizer) getTokens() []Token {
	return t.tokens
}
