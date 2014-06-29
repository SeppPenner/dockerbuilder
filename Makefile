all: build

deps:
	go get -t -v .

test:
	go test ./...

build:
	go build
