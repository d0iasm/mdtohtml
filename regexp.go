package main

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	heading, _    = regexp.Compile("(^#{1,6}) (.+)")
	headingIn, _  = regexp.Compile("^[^#]+(#{1,6}) (.+)")
	list, _       = regexp.Compile("^( *)- (.+)")
	link, _       = regexp.Compile(".*(\\[.+\\])(\\(.+\\)).*")
	whitespace, _ = regexp.Compile("^( +)(.*)")
)

type Type int

// only block elements
const (
	Newline Type = iota
	P
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

func convert(line string) Line {
	// newline
	if line == "\n" || len(line) == 0 {
		return Line{Newline, " ", 0}
	}

	// inline elements are replaced with HTML in this function.
	for link.MatchString(line) { // links <a href="url">link text</a>
		//line[loc[2]:loc[3]]: link text
		//line[loc[4]:loc[5]]: url
		loc := link.FindStringSubmatchIndex(line)

		text := line[loc[2]+1 : loc[3]-1]
		url := line[loc[4]+1 : loc[5]-1]

		litag := "<a href=\"" + url + "\">" + text + "</a>"
		line = line[:loc[2]] + litag + line[loc[5]:]
	}

	// heading existing in another component
	if headingIn.MatchString(line) {
		//line[loc[2]:loc[3]]: #, ##, ..., or ######
		//line[loc[4]:loc[5]]: title
		loc := headingIn.FindStringSubmatchIndex(line)

		n := loc[3] - loc[2] // heading number
		htag := "<h" + strconv.Itoa(n) + ">" + line[loc[4]:loc[5]] + "</h" + strconv.Itoa(n) + ">"
		line = line[:loc[2]] + htag
	}

	// break at the end of line
	if len(line) > 2 && line[len(line)-2:] == "  " {
		line = line[:len(line)-2] + "<br>"
	}

	// block elements will be replaced with HTML in the generate().
	if list.MatchString(line) {
		//line[loc[2]:loc[3]]: white spaces before "-"
		//line[loc[4]:loc[5]]: list content
		loc := list.FindStringSubmatchIndex(line)
		dep := loc[3] / 2
		return Line{Li, line[loc[4]:loc[5]], dep}
	}

	if heading.MatchString(line) {
		//line[loc[2]:loc[3]]: #, ##, ..., or ######
		//line[loc[4]:loc[5]]: title
		loc := heading.FindStringSubmatchIndex(line)
		n := loc[3]
		return Line{ntoh(n), line[loc[4]:loc[5]], 0}
	}

	// replace white spaces with a white space at the start of a line
	if whitespace.MatchString(line) {
		//line[loc[2]:loc[3]]: whitespace
		//line[loc[4]:loc[5]]: content
		loc := whitespace.FindStringSubmatchIndex(line)
		line = " " + line[loc[4]:loc[5]]
	}
	return Line{P, line, 0}
}
