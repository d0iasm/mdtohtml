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

func (t *Tokenizer) text() string {
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
	tokens := []Token{}
	for t.i < len(t.chars) {
		switch string(t.chars[t.i]) {
		case "#":
			fmt.Println("Found #", t.i)
			count := 0
			for string(t.chars[t.i]) != " " {
				t.i++
				count++
				fmt.Println(t.i, count)
				if count > 3 {
					t.i -= count
					tokens = append(tokens, Token{P, t.text()})
					break
				}
			}

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
			tokens = append(tokens, Token{RAWTEXT, t.text()})
		case "\n":
			t.i++
		default:
			tokens = append(tokens, Token{P, t.text()})
			t.i++
		}
	}
	return tokens
}
