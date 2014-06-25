PKGS := \
config \
handler \
helpers \
repository \
worker \
workspace
PKGS := $(addprefix github.com/brocaar/dockerbuilder/,$(PKGS))

all: build

deps:
	go get -t -v .

test:
	go test -v $(PKGS)

build:
	go build
