package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
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

	lines := make([]Line, 0)
	for {
		line, err := reader.ReadString('\n')
		if err != nil { // io.EOF
			break
		}
		// Remove newline at the end of line
		if len(line) > 1 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}
		lines = append(lines, convert(line))
	}

	writer.WriteString("<body>")
	writer.WriteString(generate(lines))
	writer.WriteString("</body>")
	writer.Flush()
}
