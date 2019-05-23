# mdtohtml
Mdtohtml is a HTML generator from a Markdown file. This is implemented in Go. The unique feature of this tool is that default result file includes CSS, but you can skip to include CSS by passing `-nocss` parameter.

The syntax of Markdown follows [CommonMark](https://commonmark.org/) which version is [0.29 (2019-04-06)](https://spec.commonmark.org/).

## Usage
```
$ make mdtohtml & ./mdtohtml <markdown-filename>

// You can avoid to generate css file with -nocss flag in order to customize style.
$ go build mdtohtml & ./mdtohtml <markdown-filename> -nocss
```

Current support notations (2019-05-23)
- # This is an H1
- ## This is an H2
- ### This is an H3
- #### This is an H4
- ##### This is an H5
- ###### This is an H6
- List
- Nested sublist
- Paragraph text.
- \[text\]\(url\): \<a href="url"\>text\</a\>

## EBNF
Entended Backus-Naur form for Markdown grammer.
```
Document = { Block }, EOF ;
Block = Paragraph | Headings | List ;
Paragraph = String, { String }, Newline ;
Lists = List, (List | Lists)* ;
List = ((" ")*, "-", " ", Paragraph) | Lists ;
Newline = "\n" ;
H1 = "#", String ;
H2 = "#" * 2, String ;
H3 = "#" * 3, String ;
H4 = "#" * 4, String ;
H5 = "#" * 5, String ;
H6 = "#" * 6, String ;
```

