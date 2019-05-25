#!/bin/bash
try() {
  expected="$1"
  input="$2"
  echo "$input" > test.md

  ./mdtohtml test.md -nocss
  actual="$(cat test.html)"

  if [ "$actual" = "$expected" ]; then
    echo "$input => $actual"
  else
    echo "input =>"
    echo "$input"
    echo "$expected expected, but got $actual"
    exit 1
  fi
}

echo "========== Basic Grammer =========="
echo "========== Paragrah =========="
try "<p>A paragrah.</p>" "A paragrah."

echo "========== Heading =========="
try "<h1>h1</h1>" "# h1"
try "<h2>h2</h2>" "## h2"
try "<h3>h3</h3>" "### h3"
try "<h4>h4</h4>" "#### h4"
try "<h5>h5</h5>" "##### h5"
try "<h6>h6</h6>" "###### h6"
try "<p>####### h7</p>" "####### h7"
try "<p>###dummyh3</p>" "###dummyh3"

echo "========== List =========="
try "<ul><li>list1</li></ul>" "- list1"
try "<ul><li>list1</li><li>list2</li></ul>" $'- list1\n- list2'
try "<ul><li>list1<ul><li>sublist1</li></ul></li></ul>" $'- list1\n  - sublist1'
try "<ul><li>list1<ul><li>sublist1<ul><li>subsublist1</li></ul></li></ul></li></ul>" $'- list1\n  - sublist1\n    - subsublist1'
try "<ul><li>list1<ul><li>sublist1</li></ul></li><li>list2</li></ul>" $'- list1\n  - sublist1\n- list2'
try "<ul><li>a<ul><li>aa<ul><li>aaa</li></ul></li></ul></li><li>b</li></ul>" $'- a\n  - aa\n    - aaa\n- b'
try "<ul><li>a<ul><li>aa<ul><li>aaa</li></ul></li><li>bb</li></ul></li></ul>" $'- a\n  - aa\n    - aaa\n  - bb'
try "<p>-dummylist1</p>" "-dummylist1"
# Currently, this test fails because 'c' is interpreted as a start of a new list, which means <ul>.
# try "<ul><li>a -b</li><li>c</li></ul>" $'- a\n  -b\n- c'

echo "========== Link =========="
try "<a href=\"http://example.com\">link</a>" "[link](http://example.com)"
try "<a href=\"http://example.com\">link(2)</a>" "[link(2)](http://example.com)"
try "<p>[dummylink] (http://example.com)</p>" "[dummylink] (http://example.com)"


echo "========== Combination =========="
echo "========== Heading with Inline Element =========="
try "<h1><a href=\"http://example.com\">link</a></h1>" "# [link](http://example.com)"
try "<h1>- dummylist</h1>" "# - dummylist"

echo "========== List with Inline Element =========="
try "<ul><li><a href=\"http://example.com\">link</a></li></ul>" "- [link](http://example.com)"
# try "<ul><li>This is <a href=\"http://example.com\">link</a> list.</li></ul>" "- This is [link](http://example.com) list."

echo "========== Heading after List =========="
try "<ul><li>list1</li></ul><h1>h1</h1>" $'- list1\n# h1'
try "<ul><li>list1</li></ul><h1>h1</h1>" $'- list1\n\n# h1'
try "<ul><li>a<ul><li>b</li></ul></li></ul><h1>h1</h1>" $'- a\n  - b\n# h1'

echo "========== Multiple lines =========="
try "<h1>h1</h1><p>text</p>" $'# h1\ntext'

rm test.md
rm test.html

echo OK
