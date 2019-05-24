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
	case UL:
		return "<ul>" + html + "</ul>"
	case LIST:
		return "<li>" + html + "</li>"
	case LINK:
		return "<a href=\"" + html + "\">" + node.val + "</a>"
	case URI:
		return node.val
	case P:
		return "<p>" + html + "</p>"
	default:
		return html
	}
}
