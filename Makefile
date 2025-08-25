build:
	@go build -o bin/goBackendRest cmd/main.go


test:
	@go test -V ./..


run: build
	@./bin/goBackendRest 

dev:
	@go run cmd/main.go