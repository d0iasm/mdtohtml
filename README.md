# mdtohtml
Generate a html file and a css file from a markdown file by Go.

```
$ go build mdtohtml & ./mdtohtml <markdown-filename>

// You can avoid to generate css file with -nocss flag in order to customize style.
$ go build mdtohtml & ./mdtohtml <markdown-filename> -nocss
```

Current support notations (03/12/2019)
- #: h1
- ##: h2
- ###: h3
- \[text\]\(url\): \<a href="url"\>text\</a\>
- inline/block text

## BNF(Backus-Naur form)
```
<html> ::= <html> | <h1> | <h2> | <h3> | <ul> | <p> | <br>

<h1> ::= "# " <rawtext> <newline>
<h2> ::= "## " <rawtext> <newline>
<h3> ::= "### " <rawtext> <newline>
<ul> ::= <li>
<li> ::= <li> <li> | <li> <ul>
<li> ::= "- " <rawtext> <newline> | "* " <rawtext> <newline>
<p> ::= <rawtext> <newline>
<link> ::= <url> "(" <rawtext> ")"
<url> ::= "[" <rawtext> "]"
<br> ::= <newline>
<i> ::= "*" <rawtext> "*"| "_" <rawtext> "_" 
<b> ::= "**" <rawtext> "**" | "__" <rawtext> "__" 
<text> :: = <text> | <b> | <i> | <link> | <rawtext> 

<rawtext> ::= <character> | <rawtext> 
<character> ::= <letter> | <digit> | <symbol>
<letter> ::= "A" | "B" | "C" | "D" | "E" | "F" | "G" | "H" | "I" | "J" | "K" | "L" | "M" | "N" | "O" | "P" | "Q" | "R" | "S" | "T" | "U" | "V" | "W" | "X" | "Y" | "Z" | "a" | "b" | "c" | "d" | "e" | "f" | "g" | "h" | "i" | "j" | "k" | "l" | "m" | "n" | "o" | "p" | "q" | "r" | "s" | "t" | "u" | "v" | "w" | "x" | "y" | "z"
<digit> ::= "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
<symbol> ::= "#" | "*" | "-" | "[" | "]" | "(" | ")"  
<newline> ::= "\n" | "\r\n"
```
