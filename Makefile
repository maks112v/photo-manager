
.PHONY: build
build:
	rm -rf bin
	go build -o bin/photomanager ./main.go

.PHONY: generate
generate:
	go generate ./...