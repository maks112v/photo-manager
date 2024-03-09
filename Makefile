
.PHONY: build
build:
	go build -o bin/photomanager ./main.go

.PHONY: generate
generate:
	go generate ./...