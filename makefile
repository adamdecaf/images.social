.PHONY: build test vet

build: vet
	go build .

vet:
	go tool vet .

test: build
	go test -v .
