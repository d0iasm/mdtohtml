package main

import (
	"fmt"
)

type Node struct {
	ty       Type
	children []Node
	val      string
}

type Parser struct {
	i      int
	tokens []Token
}

func appendChild(parent *Node, child Node) {
	parent.children = append(parent.children, child)
}

func (p *Parser) heading(level Type) Node {
	t := p.tokens[p.i]
	n := Node{level, []Node{}, ""}
	p.i++
	t = p.tokens[p.i]
	if t.ty == RAWTEXT {
		appendChild(&n, p.rawtext(t.val))
	}
	return n
}

func (p *Parser) p() Node {
	t := p.tokens[p.i]
	n := Node{P, []Node{}, ""}
	appendChild(&n, p.rawtext(t.val))
	return n
}

func (p *Parser) rawtext(s string) Node {
	return Node{RAWTEXT, []Node{}, s}
}

func (p *Parser) body() Node {
	root := Node{BODY, []Node{}, ""}
	for p.i < len(p.tokens) {
		t := p.tokens[p.i]
		fmt.Println("Called body", t, p.i)
		switch t.ty {
		case H1:
			appendChild(&root, p.heading(H1))
		case H2:
			appendChild(&root, p.heading(H2))
		case H3:
			appendChild(&root, p.heading(H3))
		case P:
			appendChild(&root, p.p())
		}
		p.i++
	}
	return root
}