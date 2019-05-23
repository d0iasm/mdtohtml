package main

import (
	"fmt"
	"strings"
)

type Node struct {
	ty       Type
	children []Node
        parent Node
	val      string
}

type Parser struct {
	i      int
	tokens []Token
}

func appendChild(parent *Node, child Node) {
	parent.children = append(parent.children, child)
}

func (p *Parser) heading() Node {
	t := p.tokens[p.i]
	n := Node{HEADING, []Node{}, t.val}
	p.i++
	t = p.tokens[p.i]
	if t.ty != RAWTEXT {
		panic("Token next to a heading should be a raw text.")
	}
	appendChild(&n, p.rawtext(t.val))
	return n
}

func (p *Parser) headingWith(level Type) Node {
	t := p.tokens[p.i]
	n := Node{level, []Node{}, t.val}
	p.i++
	t = p.tokens[p.i]
	if t.ty != RAWTEXT {
		panic("Token next to a heading should be a raw text.")
	}
	appendChild(&n, p.rawtext(t.val))
	return n
}

/**
func (p *Parser) ul(parent Node) Node {
	fmt.Println("UL: ", p.i, p.tokens[p.i])
	n := Node{UL, []Node{}, ""}
	p.i++
	t := p.tokens[p.i]
	for t.ty == UL || t.ty == LIST {
		if p.i++; p.i >= len(p.tokens) {
			return n
		}

		switch t.ty {
		case UL:
			appendChild(&n, p.ul(n))
		case LIST:
			appendChild(&n, p.list(n))
		}
		if p.i++; p.i >= len(p.tokens) {
			p.i--
			return n
		}
		t = p.tokens[p.i]
	}
                return n
}

func (p *Parser) list(parent Node) Node {
	fmt.Println("LIST: ", p.i, p.tokens[p.i])
	n := Node{LIST, []Node{}, ""}
		if p.i++; p.i >= len(p.tokens) {
			return n
		}

                t := p.tokens[p.i]
		switch t.ty {
		case UL:
			appendChild(&n, p.ul(n))
		case UL_END:
                        appendChild(&parent, n)
			return n
		case LINK:
			appendChild(&n, p.link())
		case RAWTEXT:
			appendChild(&n, p.rawtext(t.val))
		}
                return n
}
*/

func (p *Parser) ul() Node {
	fmt.Println("UL: ", p.i, p.tokens[p.i])
	n := Node{UL, []Node{}, ""}
	p.i++
	t := p.tokens[p.i]
	for t.ty == UL || t.ty == LIST {
		fmt.Println("UL: ", p.i, p.tokens[p.i], n)
		switch t.ty {
		case UL:
			appendChild(&n, p.ul())
		case LIST:
			appendChild(&n, p.list(n))
		}

		if p.i++; p.i >= len(p.tokens) {
			p.i--
			return n
		}
		t = p.tokens[p.i]
	}
	p.i--
	return n
}

func (p *Parser) list(parent Node) Node {
	fmt.Println("LIST: ", p.i, p.tokens[p.i])
	t := p.tokens[p.i]
	n := Node{LIST, []Node{}, t.val}
	p.i++
	t = p.tokens[p.i]
	for t.ty == UL || t.ty == UL_END || t.ty == LINK || t.ty == RAWTEXT {
		fmt.Println("LIST: ", p.i, p.tokens[p.i], n)
		switch t.ty {
		case UL:
			appendChild(&n, p.ul())
		case UL_END:
			appendChild(&parent, n)
			return n
		case LINK:
			appendChild(&n, p.link())
		case RAWTEXT:
			appendChild(&n, p.rawtext(t.val))
		}

		if p.i++; p.i >= len(p.tokens) {
			p.i--
			return n
		}
		t = p.tokens[p.i]
	}
	p.i--
	return n
}

func (p *Parser) link() Node {
	t := p.tokens[p.i]
	n := Node{LINK, []Node{}, t.val}

	p.i++
	t = p.tokens[p.i]
	if t.ty != RAWTEXT {
		panic("Token next to LINK should be rawtext.")
	}
	appendChild(&n, p.rawtext(t.val))
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
		//fmt.Println("Called body", t, p.i)
		switch t.ty {
		case HEADING:
			appendChild(&root, p.heading())
		case UL:
			fmt.Println("root UL: ", p.i, p.tokens[p.i])
			appendChild(&root, p.ul())
		case LINK:
			appendChild(&root, p.link())
		case P:
			appendChild(&root, p.p())
		case H1:
			appendChild(&root, p.headingWith(H1))
		case H2:
			appendChild(&root, p.headingWith(H2))
		case H3:
			appendChild(&root, p.headingWith(H3))
		default:
			appendChild(&root, p.rawtext(t.val))
		}
		p.i++
	}
	return root
}

func debugTree(root Node, dep int) {
	spaces := strings.Repeat(" ", dep*2)
	if dep == 0 {
		fmt.Println(spaces+"Node (type, val, dep):", root.ty, root.val, dep)
	}
	fmt.Println(spaces + "|")
	for _, c := range root.children {
		fmt.Println(spaces+"  - Node (type, val, dep):", c.ty, c.val, dep+1)
		debugTree(c, dep+1)
	}
}
