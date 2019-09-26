package main

import (
	"fmt"
	"strings"
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

func (p *Parser) heading() Node {
	t := p.tokens[p.i]
	n := Node{HEADING, []Node{}, t.val}
	p.i++
	t = p.tokens[p.i]
	if t.ty != RAWTEXT && t.ty != LINK {
		panic("The token next to heading should be inline element.")
	}

	switch t.ty {
	case LINK:
		appendChild(&n, p.link())
	case RAWTEXT:
		appendChild(&n, p.rawtext())
	}
	return n
}

func (p *Parser) ul(dep int) Node {
	n := Node{UL, []Node{}, ""}
	for p.tokens[p.i].ty != EOF {
		p.i++
		switch p.tokens[p.i].ty {
		case UL:
			appendChild(&n, p.ul(dep+1))
		case LIST:
			if p.tokens[p.i].dep < dep {
				p.i--
				return n
			}
			appendChild(&n, p.list(dep))
		default:
			p.i--
			return n
		}
	}
	return n
}

func (p *Parser) list(dep int) Node {
	n := Node{LIST, []Node{}, p.tokens[p.i].val}
	for p.tokens[p.i].ty != EOF {
		p.i++
		switch p.tokens[p.i].ty {
		case UL:
			appendChild(&n, p.ul(dep+1))
		case LINK:
			appendChild(&n, p.link())
		case HEADING:
			appendChild(&n, p.heading())
		case RAWTEXT:
			appendChild(&n, p.rawtext())
		default:
			p.i--
			return n
		}
	}
	return n
}

func (p *Parser) link() Node {
	t := p.tokens[p.i]
	n := Node{LINK, []Node{}, t.val}

	p.i++
	t = p.tokens[p.i]
	if t.ty != URI {
		p.i--
		return p.rawtext()
	}
	appendChild(&n, Node{URI, []Node{}, t.val})
	return n
}

func (p *Parser) p() Node {
	n := Node{P, []Node{}, ""}
	appendChild(&n, p.rawtext())
	return n
}

func (p *Parser) rawtext() Node {
	return Node{RAWTEXT, []Node{}, p.tokens[p.i].val}
}

func (p *Parser) body() Node {
	root := Node{BODY, []Node{}, ""}
	for p.i < len(p.tokens) {
		t := p.tokens[p.i]
		switch t.ty {
		case HEADING:
			appendChild(&root, p.heading())
		case UL:
			appendChild(&root, p.ul(0))
		case LINK:
			appendChild(&root, p.link())
		case P:
			appendChild(&root, p.p())
		case EOF:
			return root
		default:
			appendChild(&root, p.rawtext())
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
