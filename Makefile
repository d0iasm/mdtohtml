files=mdtohtml.go tokenizer.go parser.go generator.go css.go

mdtohtml: $(files) 
	go build $(files) 

clean:
	rm mdtohtml
