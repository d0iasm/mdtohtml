package main

func css() string {
	return `
<link href="https://fonts.googleapis.com/css?family=Abril+Fatface|Lora|Noto+Serif|Noto+Serif+JP|Trirong:700" rel="stylesheet">
<style>
body {
  font-family: 'Noto Sans', sans-serif;
  font-size: 18px;
  max-width: 740px;
  margin: 30px auto 0 auto;
}
h1, h2, h3, h4, h5, h6 {
  border-bottom: 3px black solid;
  margin-bottom: 0;
}
</style>
`
}
