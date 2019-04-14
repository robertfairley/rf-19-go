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
build-deploy:
	@make clean
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GOOS=$(DEPLOY_OS) go build -o bin/$(GONAME) ./src/$(GOFILES)
test:
	@echo "Running tests..."
	@echo "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
	@GOPATH=$(GOPATH) go test ./src && echo "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\nPASSED"
get:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get ./src
clean:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean
	@rm -rf ./bin ./pkg
reset:
	@rm -rf ./vendor ./bin ./pkg
docs:
	@godoc -http=:5005
deploy:
	@./scripts/deploy.sh

.PHONY: build build-deploy get run start clean docs deploy
