package main

import (
	"fmt"
	"strings"
)

func test(expect string, input string) {
	lines := make([]Line, 0)
	for _, in := range strings.Split(input, "\n") {
		lines = append(lines, convert(in))
	}
	html := generate(lines)

	if html == expect {
		fmt.Println(input + " => " + expect)
	} else {
		panic(input + " => " + expect + " but got " + html)
	}
}

func main() {
	fmt.Println("\n----- paragrah -----")
	test("<p>a paragraph</p>", "a paragraph")
	test("<p>a paragraph<br>hogehoge</p>", "a paragraph  \nhogehoge")

	fmt.Println("\n----- heading -----")
	test("<h1>h1</h1>", "# h1")
	test("<h2>h2</h2>", "## h2")
	test("<h3>h3</h3>", "### h3")
	test("<h4>h4</h4>", "#### h4")
	test("<h5>h5</h5>", "##### h5")
	test("<h6>h6</h6>", "###### h6")
	test("<p>####### h7</p>", "####### h7")
	test("<p>###dummyh3</p>", "###dummyh3")

	fmt.Println("\n----- list -----")
	test("<ul><li>list1</li></ul>", "- list1")
	test("<ul><li>list1</li><li>list2</li></ul>", "- list1\n- list2")
	// Sublist is not a standard syntax.
	// It should be <li>list1<ul><li>sublist1</li></ul></li></ul>
	// but now got <li>list1</li><ul><li>sublist1</li></ul></ul>
	test("<ul><li>list1</li><ul><li>sublist1</li></ul></ul>", "- list1\n  - sublist1")
	test("<ul><li>list1</li><ul><li>sublist1</li><ul><li>subsublist1</li></ul></ul></ul>", "- list1\n  - sublist1\n    - subsublist1")
	test("<ul><li>list1</li><ul><li>sublist1</li></ul><li>list2</li></ul>", "- list1\n  - sublist1\n- list2")
	test("<ul><li>a</li><ul><li>aa</li><ul><li>aaa</li></ul></ul><li>b</li></ul>", "- a\n  - aa\n    - aaa\n- b")
	test("<ul><li>a</li><ul><li>aa</li><ul><li>aaa</li></ul><li>bb</li></ul></ul>", "- a\n  - aa\n    - aaa\n  - bb")
	test("<ul><li><h1>h1</h1></li></ul>", "- # h1")

	//Currently, this test fails because "c" is interpreted as a start of a new list, which means <ul>.
	//test("<ul><li>a -b</li><li>c</li></ul>", "- a\n  -b\n- c")

	fmt.Println("\n----- link -----")
	test("<p><a href=\"http://example.com\">link</a></p>", "[link](http://example.com)")
	test("<p><a href=\"http://example.com\">link(2)</a></p>", "[link(2)](http://example.com)")
	test("<p>inline text<a href=\"http://example.com\">link</a>.</p>", "inline text[link](http://example.com).")
	test("<p>[dummylink] (http://example.com)</p>", "[dummylink] (http://example.com)")

	fmt.Println("\n----- heading with inline elements -----")
	test("<h1><a href=\"http://example.com\">link</a></h1>", "# [link](http://example.com)")
	test("<h1>- dummylist</h1>", "# - dummylist")

	fmt.Println("\n----- list with inline elements -----")
	test("<ul><li><a href=\"http://example.com\">link</a></li></ul>", "- [link](http://example.com)")
	test("<ul><li>This is <a href=\"http://example.com\">link</a> list.</li></ul>", "- This is [link](http://example.com) list.")

	fmt.Println("\n----- heading after a list -----")
	test("<ul><li>list1</li></ul><h1>h1</h1>", "- list1\n# h1")
	test("<ul><li>list1</li></ul><h1>h1</h1>", "- list1\n\n# h1")
	test("<ul><li>a</li><ul><li>b</li></ul></ul><h1>h1</h1>", "- a\n  - b\n# h1")

	fmt.Println("\n----- multiple lines -----")
	test("<h1>h1</h1><p>text</p>", "# h1\ntext")

	fmt.Println("OK")
}
