all: build

deps:
	go get -t -v .

test:
	go test -v ./...

build:
	go build
