.PHONY: build tidy fmt

build:
	go build -o download-cwlogs main.go

tidy:
	go mod tidy

fmt:
	go fmt ./...
