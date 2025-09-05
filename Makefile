
setup: 
	@go mod tidy
	@go mod download
	@mkdir -p bin
	@docker compose -f compose.yaml up -d
	@go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	
	@make migrate-up



clean-setup:
	@docker compose -f compose.yaml down -v --rmi all --remove-orphans
	@rm -rf bin
	@make setup

build:
	@go build -o bin/goBackendRest cmd/main.go


test:
	@go test -v ./...


run: build
	@./bin/goBackendRest 

dev:
	@go run cmd/main.go

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down