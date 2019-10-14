# mdtohtml
Mdtohtml is a HTML generator from a Markdown file. This is implemented in Go. The unique feature of this tool is that default result file includes CSS, but you can skip to include CSS by passing `-nocss` parameter.

The syntax of Markdown follows [CommonMark](https://commonmark.org/) which version is [0.29 (2019-04-06)](https://spec.commonmark.org/). They have playground, [commonmark.js dingus](https://spec.commonmark.org/dingus/), for CommonMark grammar.

## Usage
```
$ make mdtohtml & ./mdtohtml <markdown-filename>

// You can avoid to generate css file with -nocss flag in order to customize style.
$ make mdtohtml & ./mdtohtml <markdown-filename> -nocss
```

## Current support notations (2019-10-14)
- Paragraph
- Headings
- List
- Nested list
- Link
- Emphasis
- Strong

## EBNF
Entended Backus-Naur form for Markdown grammer.
```
Document = { Block }, EOF ;
Block = Paragraph | Headings | Lists ;
Inline = Link | Rawtext | Emphasis | Strong ;
Paragraph = String, { String }, Newline ;
Headings = H1 | H2 | H3 | H4 | H5 | H6 ;
H1 = "#", Inline ;
H2 = "#" * 2, Inline ;
H3 = "#" * 3, Inline ;
H4 = "#" * 4, Inline ;
H5 = "#" * 5, Inline ;
H6 = "#" * 6, Inline ;
Lists = List, ( List | Lists )* ;
List = ( (" ")*, "-", " ", Inline ) | Lists ;
Emphasis = "*" Rawtext "*" | "_" Rawtext "_" ;
Strong = "**" Rawtext "**" | "__" Rawtext "__" ;
String = { Character } ;
Newline = "\n" ;
```

## Test
Supports test written in Go. You can test this tool just by `$ make test`.
