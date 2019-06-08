package main

func css() string {
	return `
<style>
body {
  font-size: 18px;
  max-width: 680px;
  margin: 30px auto 0 auto;
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
p {
  margin: 0.5em auto;
}
ul {
  margin: 0.5em auto 1em auto;
}
</style>
`
}
