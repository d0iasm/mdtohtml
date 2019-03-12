files=mdtohtml.go tokenizer.go parser.go generator.go

mdtohtml: mdtohtml.go
	go build $(files) 

clean:
	rm mdtohtml
