package main

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	heading, _   = regexp.Compile("(^#{1,6}) (.+)")
	headingIn, _ = regexp.Compile("^[^#]+(#{1,6}) (.+)")
	list, _      = regexp.Compile("^( *)- (.+)")
	link, _      = regexp.Compile(".*(\\[.+\\])(\\(.+\\)).*")
)

type Type int

// only block elements
const (
	P Type = iota
	H1
	H2
	H3
	H4
	H5
	H6
	Li
)

type Line struct {
	ty  Type
	val string
	dep int
}

func ntoh(n int) Type {
	switch n {
	case 1:
		return H1
	case 2:
		return H2
	case 3:
		return H3
	case 4:
		return H4
	case 5:
		return H5
	case 6:
		return H6
	default:
		panic(fmt.Sprintf("a heading should be in the range of 1 to 6, but got %d", n))
	}
}

func hton(ty Type) int {
	switch ty {
	case H1:
		return 1
	case H2:
		return 2
	case H3:
		return 3
	case H4:
		return 4
	case H5:
		return 5
	case H6:
		return 6
	default:
		panic(fmt.Sprintf("a heading should be in the range of 1 to 6, but got %d", ty))
	}
}

func transpile(line []byte) Line {
	// inline elements are replaced with HTML in this function.
	for link.Match(line) { // links <a href="url">link text</a>
		//line[loc[2]:loc[3]]: link text
		//line[loc[4]:loc[5]]: url
		loc := link.FindSubmatchIndex(line)

		text := make([]byte, loc[3]-loc[2]-2)
		copy(text, line[loc[2]+1:loc[3]-1])
		url := make([]byte, loc[5]-loc[4]-2)
		copy(url, line[loc[4]+1:loc[5]-1])

		newli := make([]byte, 0)
		newli = append(newli, line[:loc[2]]...)
		newli = append(newli, []byte("<a href=\"")...)
		newli = append(newli, url...)
		newli = append(newli, []byte("\">")...)
		newli = append(newli, text...)
		newli = append(newli, []byte("</a>")...)
		newli = append(newli, line[loc[5]:]...)

		// extend the length
		line = make([]byte, len(newli))
		copy(line, newli)
	}

	if headingIn.Match(line) {
		//line[loc[2]:loc[3]]: #, ##, ..., or ######
		//line[loc[4]:loc[5]]: title
		loc := headingIn.FindSubmatchIndex(line)
		n := loc[3] - loc[2]
		h := "h" + strconv.Itoa(n)
		newh := "<" + h + ">" + string(line[loc[4]:loc[5]]) + "</" + h + ">"
		newh = string(line[:loc[2]]) + newh
		line = make([]byte, len([]byte(newh)))
		copy(line, []byte(newh))
	}

	// block elements will be replaced with HTML in the generate().
	if list.Match(line) {
		//line[loc[2]:loc[3]]: white spaces before "-"
		//line[loc[4]:loc[5]]: list content
		loc := list.FindSubmatchIndex(line)
		dep := loc[3] / 2
		return Line{Li, string(line[loc[4]:loc[5]]), dep}
	}

	if heading.Match(line) {
		//line[loc[2]:loc[3]]: #, ##, ..., or ######
		//line[loc[4]:loc[5]]: title
		loc := heading.FindSubmatchIndex(line)
		n := loc[3]
		return Line{ntoh(n), string(line[loc[4]:loc[5]]), 0}
	}
	return Line{P, string(line), 0}
}
