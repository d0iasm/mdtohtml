package main

import (
	"fmt"
)

type Token struct {
	ty  Type
	val string
}

type Tokenizer struct {
	i     int
	chars []rune
}

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

func (t *Tokenizer) tokenize() []Token {
	fmt.Println(string(t.chars))
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
