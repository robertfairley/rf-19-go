build:
	make clean
	go build -o bin/server src/main.go
clean:
	rm -rf ./bin
setup:
	go get -ugopkg.in/russross/blackfriday.v2
