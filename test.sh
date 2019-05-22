#!/bin/bash
try() {
  expected="$1"
  input="$2"
  echo $input > test.md

  ./mdtohtml test.md -nocss
  actual="$(cat test.html)"

  if [ "$actual" = "$expected" ]; then
    echo "$input => $actual"
  else
    echo "$expected expected, but got $actual"
    exit 1
  fi
}

echo "========== Basic Grammer =========="
try "<h1>h1</h1>" "# h1"
try "<h2>h2</h2>" "## h2"
try "<h3>h3</h3>" "### h3"
try "<h4>h4</h4>" "#### h4"
try "<h5>h5</h5>" "##### h5"
try "<h6>h6</h6>" "###### h6"
try "<p>####### h7</p>" "####### h7"
try "<p>###dummyh3</p>" "###dummyh3"

rm test.md
rm test.html

echo OK
