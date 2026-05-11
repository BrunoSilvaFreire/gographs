.PHONY: all build test clean

all: build

build:
	go build ./...

test:
	go test -v ./...

clean:
	go clean
	rm -f *.deb *.zip
