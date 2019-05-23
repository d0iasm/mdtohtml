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
}

type Tokenizer struct {
	s      *bufio.Scanner
	buf    *bytes.Buffer
	tokens []Token
}

/**
func (t *Tokenizer) checkUntilEnd(target string) int {
	i := 0
	for string(t.chars[t.i]) != "\n" {
		if string(t.chars[t.i]) == target {
			return i
		}
		t.i++
		i++
	}
	return -1
}

func (t *Tokenizer) rawtext() string {
	start := t.i
	for string(t.chars[t.i]) != "\n" {
		t.i++
		if t.i >= len(t.chars) {
			break
		}
	}
	return string(t.chars[start:t.i])
}
*/

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
}

func (t *Tokenizer) list(nest int, sym string) {
	t.tokens = append(t.tokens, Token{LIST, sym})
	t.tokens = append(t.tokens, Token{RAWTEXT, t.stringLiteral()})
	if t.s.Scan() {
		t.buf.WriteString(t.s.Text())
		switch t.s.Text() {
		case sym:
			fmt.Println("nest: ", nest)
			fmt.Println("Called sym!!!", t.buf.String(), strings.Count(t.buf.String(), " "))

			if strings.Count(t.buf.String(), " ") < (nest+1)*2 {
				return
			}
			if t.consume(" ") {
				t.list(nest+1, sym)
			}
		case " ":
			n := t.count(" ")
			fmt.Println("nest: ", nest)
			fmt.Println("Called space!!!", t.buf.String(), strings.Count(t.buf.String(), " "), n)
			if n == (nest+1)*2 && t.consume(" ") {
				t.buf.Reset()
				t.ul(nest+1, sym)
			}
		default:
			return
		}
	}

}

func (t *Tokenizer) tokenize() {
	isHead := true
	//nest := 0

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
			//nest = 0
		case "-":
			sym := t.s.Text()
			if t.consume(" ") {
				t.ul(0, sym)
			}
			t.buf.WriteString(sym)
			t.buf.WriteString(t.s.Text())
			/**
			fmt.Println(strings.Count(t.buf.String(), " "), nest)
			if strings.Count(t.buf.String(), " ") < (nest+1)*2 {

				t.tokens = append(t.tokens, Token{UL_END, ""})
				nest = 0
				break
			}
			if strings.Count(t.buf.String(), " ") == (nest+1)*2 {
				t.buf.Reset()
				nest++
			}

			if !isHead && nest == 0 {
				t.buf.WriteString(t.s.Text())
				break
			}

			sym := t.s.Text()
			if t.s.Scan() && t.checkSpace() {
				if nest > 0 || len(tokens) < 2 || tokens[len(tokens)-2].ty != LIST {
					tokens = append(tokens, Token{UL, sym})
				}
				tokens = append(tokens, Token{LIST, sym})
				isHead = false
				break
			}
			t.buf.WriteString(sym)
			t.buf.WriteString(t.s.Text())
			*/
		default:
			t.buf.WriteString(t.s.Text())
		}
	}
}

func (t *Tokenizer) getTokens() []Token {
	return t.tokens
}

/**
func (t *Tokenizer) tokenize_old() []Token {
	buf := []rune{}
	tokens := []Token{}
	isHead := true
	shouldInline := false

	for t.i < len(t.chars) {
		switch string(t.chars[t.i]) {
		case "#":
			if !isHead {
				buf = append(buf, t.chars[t.i])
				t.i++
				break
			}

			posNextWhitespace := t.checkUntilEnd(" ")
			if posNextWhitespace > 3 {
				buf = append(buf, t.chars[t.i-posNextWhitespace:t.i]...)
				break
			}

			if len(buf) > 0 {
				tokens = append(tokens, Token{RAWTEXT, string(buf)})
				buf = buf[:0]
			}

			count := posNextWhitespace
			if count == 1 {
				tokens = append(tokens, Token{H1, "#"})
			} else if count == 2 {
				tokens = append(tokens, Token{H2, "##"})
			} else if count == 3 {
				tokens = append(tokens, Token{H3, "###"})
			} else {
				break
			}
			t.i++
			isHead = false
			tokens = append(tokens, Token{RAWTEXT, t.rawtext()})
		case "-":
			if !isHead {
				buf = append(buf, t.chars[t.i])
				t.i++
				break
			}

			posNextWhitespace := t.checkUntilEnd(" ")
			if posNextWhitespace > 1 {
				buf = append(buf, t.chars[t.i-posNextWhitespace:t.i]...)
				break
			}

			tokens = append(tokens, Token{LIST, "-"})
			t.i++
			isHead = false
			shouldInline = true
		case "[":
			start := t.i
			posEndText := t.checkUntilEnd("]")
			posStartUrl := t.checkUntilEnd("(")
			posEndUrl := t.checkUntilEnd(")")
			if posEndText != -1 && posStartUrl == 1 && posEndUrl != -1 {
				if len(buf) > 0 {
					tokens = append(tokens, Token{RAWTEXT, string(buf)})
					buf = buf[:0]
				}
				tokens = append(tokens, Token{LINK, string(t.chars[start+posEndText+posStartUrl+1 : t.i])})
				tokens = append(tokens, Token{RAWTEXT, string(t.chars[start+1 : start+posEndText])})
				t.i++
				isHead = false
				shouldInline = true
			} else {
				t.i = start + 1
			}
		case "\n":
			if len(buf) > 0 && shouldInline {
				tokens = append(tokens, Token{RAWTEXT, string(buf)})
				buf = buf[:0]
			} else if len(buf) > 0 {
				tokens = append(tokens, Token{P, string(buf)})
				buf = buf[:0]
			}
			t.i++
			isHead = true
			shouldInline = false
		default:
			buf = append(buf, t.chars[t.i])
			t.i++
			isHead = false
		}
	}
	return tokens
}
*/
