CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test ! -d src/github.com/straup/go-image-tools; then mkdir -p src/github.com/straup/go-image-tools; fi
	cp -r util src/github.com/straup/go-image-tools/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/iand/salience"
	@GOPATH=$(GOPATH) go get -u "github.com/jung-kurt/gofpdf"
	@GOPATH=$(GOPATH) go get -u "github.com/MaxHalford/halfgone"
	@GOPATH=$(GOPATH) go get -u "github.com/nfnt/resize/"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/crop cmd/crop.go
	@GOPATH=$(GOPATH) go build -o bin/halftone cmd/halftone.go
	@GOPATH=$(GOPATH) go build -o bin/picturebook cmd/picturebook.go