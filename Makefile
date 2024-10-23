.PHONY: build tidy fmt

build:
	go build -o download-cwlogs ./cmd/download-cwlogs

tidy:
	go mod tidy

fmt:
	go fmt ./...
