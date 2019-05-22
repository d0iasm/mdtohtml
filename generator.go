package main

import (
	"fmt"
	"strings"
)

func generate(node Node) string {
	//fmt.Println("Current Node", node.ty, node.val)
	if node.ty == RAWTEXT {
		return node.val
	}

	html := ""
	for _, c := range node.children {
		html += generate(c)
	}

	switch node.ty {
	case HEADING:
		n := strings.Count(node.val, "#")
		return fmt.Sprintf("<h%d>%s</h%d>", n, html, n)
	case H1:
		return "<h1>" + html + "</h1>"
	case H2:
		return "<h2>" + html + "</h2>"
	case H3:
		return "<h3>" + html + "</h3>"
	case LIST:
		return "<li>" + html + "</li>"
	case LINK:
		return "<a href=" + node.val + ">" + html + "</a>"
	case P:
		return "<p>" + html + "</p>"
	default:
		return html
	}
}
