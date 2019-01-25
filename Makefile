build:
	make clean
	go build -o bin/server src/main.go
clean:
	rm -rf ./bin
setup:
	export GOPATH=$(pwd)
	go get -ugopkg.in/russross/blackfriday.v2
