.PHONY: build test vet

build: vet
	go build .

vet:
	go fmt github.com/adamdecaf/images.social
	go tool vet .

test: build
	go test -v ./...
