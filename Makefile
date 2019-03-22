GOPATH=$(shell pwd)/vendor:$(shell pwd)
GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
GONAME="server"
DEPLOY_OS="linux"

run:
	@go run src/main.go
start:
	./bin/$(GONAME)
build:
	@make clean
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o bin/$(GONAME) ./src/$(GOFILES)
get:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get ./src
clean:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean
	@rm -rf ./bin ./pkg
reset:
	@rm -rf ./vendor ./bin ./pkg
docs:
	@godoc -http=:5005

.PHONY: build get run start clean docs
