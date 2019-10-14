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
	link, _    = regexp.Compile(".*(\\[.+\\])(\\(.+\\)).*")
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func compile(line []byte) string {
	// links <a href="url">link text</a>
	for link.Match(line) {
		//link[loc[2]:loc[3]]: link text
		//link[loc[4]:loc[5]]: url
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

	// headings <h1>, <h2>, <h3>, <h4>, <h5>, and <h6>
	if heading.Match(line) {
		//link[loc[2]:loc[3]]: #, ##, ..., or ######
		//link[loc[4]:loc[5]]: title
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
		fmt.Println("---------")
		fmt.Println(string(line))
		fmt.Println(html)
		fmt.Println("---------")
		writer.WriteString(html)
	}
	writer.Flush()
}
