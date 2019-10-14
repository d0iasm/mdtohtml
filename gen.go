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
			// insert <ul> for the start of lists
			if (i > 0 && lines[i-1].ty != Li) || i == 0 {
				html += "<ul>"
			}
                        // insert <ul> for the start of sublists
                        if i > 0 && lines[i-1].dep < l.dep {
				html += "<ul>"
                        }

			html += "<li>"
			html += l.val
			html += "</li>"

                        // insert </ul> for the end of sublists
			if i < len(lines)-1 && l.dep > lines[i+1].dep {
				html += "</ul>"
                        }
                        // insert </ul> for the end of sublists when a document ends with lists
                        if i == len(lines)-1 {
                          dep := l.dep
                          for dep > 0 {
                            html += "</ul>"
                            dep -= 1
                          }
                        }

			// insert </ul> for the end of lists
			if i < len(lines)-1 && lines[i+1].ty != Li || i == len(lines)-1 {
				html += "</ul>"
			}
		default:
			html += l.val
		}
	}
	return html
}
