package main

func css() string {
	return `
<link href="https://fonts.googleapis.com/css?family=Roboto&display=swap" rel="stylesheet">
<style>
body {
  font-family: 'Roboto', sans-serif;
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
h4 {
  background-color: black;
  color: white;
  padding: 0.1em 0.25em;
  margin: 0;
  width: fit-content;
}
h1 > a, h2 > a {
  text-decoration: none;
}
a:hover {
  opacity: 0.5;
}
p {
  margin: 0 auto 0.5em auto;
}
ul {
  margin: 0 auto;
}
em {


}
</style>
`
}
