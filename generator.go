package main

import (
	"fmt"
)

func generate(node Node) string {
	fmt.Println("Current Node", node.ty, node.val)
	if node.ty == RAWTEXT {
		return node.val
	}

	html := ""
	for _, c := range node.children {
		html += generate(c)
	}

	switch node.ty {
	case H1:
		return "<h1>" + html + "</h1>"
	case H2:
		return "<h2>" + html + "</h2>"
	case H3:
		return "<h3>" + html + "</h3>"
	case P:
		return "<p>" + html + "</p>"
	default:
		return html
	}
}
