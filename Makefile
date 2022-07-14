BINARY=app

download:
	@go mod download

test: 
	@go test -v -cover -covermode=atomic ./...

run: 
	@go run server.go

build:
	@go build -o ${BINARY} server.go

runall:
	@go mod download
	@go build -o ${BINARY} server.go
	./app

.PHONY: test run download build
