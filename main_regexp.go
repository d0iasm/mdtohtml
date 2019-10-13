package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type TYPE *regexp.Regexp

var (
	heading, _ = regexp.Compile("(^#{1,6}) (.+)")
	list, _    = regexp.Compile("( *)- (.+)")
	link, _    = regexp.Compile("(\\[.+\\])(\\(.+\\))")
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func compile(line []byte) string {
	// headings <h1>, <h2>, <h3>, <h4>, <h5>, and <h6>
	if heading.Match(line) {
		loc := heading.FindSubmatchIndex(line)
		n := strconv.Itoa(loc[3])
		return "<h" + n + ">" + string(line[loc[4]:loc[5]]) + "</h" + n + ">"
	}
	return string(line)
}

func main() {
	fname := os.Args[1]
	name := strings.Split(fname, ".")

	if strings.Compare(strings.ToLower(name[len(name)-1]), "md") != 0 {
		panic("input file must be a markdown file (.md)")
	}

	wfile, err := os.Create(name[0] + ".html")
	check(err)
	defer wfile.Close()
	writer := bufio.NewWriter(wfile)
	if len(os.Args) < 3 || os.Args[2] != "-nocss" {
		fmt.Fprintln(writer, css())
	}

	rfile, err := os.Open(fname)
	check(err)
	defer rfile.Close()
	reader := bufio.NewReader(rfile)

	for {
		line, _, err := reader.ReadLine()
		if err != nil { // io.EOF
			break
		}
		html := compile(line)
		writer.WriteString(html)
	}
	writer.Flush()
}
