package main

import (
	"fmt"
	"strconv"
)

func generate(lines []Line) string {
	html := ""
	for i, l := range lines {
		fmt.Println(i, l)
		switch l.ty {
		case H1, H2, H3, H4, H5, H6:
			html += "<h" + strconv.Itoa(hton(l.ty)) + ">"
			html += l.val
			html += "</h" + strconv.Itoa(hton(l.ty)) + ">"
		case Li:
			// insert <ul> because of the start of links
			if (i > 0 && lines[i-1].ty != Li) || i == 0 {
				html += "<ul>"
			}
			html += "<li>"
			html += l.val
			html += "</li>"

			// insert </ul> because of the end of links
			if i < len(lines)-1 && lines[i+1].ty != Li || i == len(lines)-1 {
				html += "</ul>"
			}
		default:
			html += l.val
		}
	}
	return html
}
