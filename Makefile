files=regexp.go gen.go css.go

mdtohtml: main.go $(files)
	go build -o mdtohtml main.go $(files)

test: test.go mdtohtml
	go build -o test test.go $(files)
	./test

clean:
	rm mdtohtml
