package main

func css() string {
	return `
<link href="https://fonts.googleapis.com/css?family=Noto+Sans+JP&display=swap" rel="stylesheet">
<style>
body {
  font-family: 'Noto Sans JP', sans-serif;
  font-size: 16px;
  max-width: 680px;
  margin: 30px auto 0 auto;
}
@media (max-width: 980px) {
  body {
    max-width: 90%;
  }
}
h1, h2, h3, h4, h5, h6 {
  margin-bottom: 0.5em;
}
h1 {
  font-size: 48px;
  text-align: center;
}
h2 {
  border-bottom: 3px black solid;
}
h1 > a, h2 > a {
  text-decoration: none;
}
a:hover {
  opacity: 0.5;
}
p, ul {
  margin: 0 auto 0.5em auto;
}
</style>
`
}
