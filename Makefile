
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
	@docker compose -f compose.yaml up -d --build
	@echo "Waiting for Docker services to be healthy..."
	@TRIES=0; \
	until docker inspect --format='{{.State.Health.Status}}' my-mysql | grep -q healthy; do \
	  echo "MySQL health status: $$(docker inspect --format='{{.State.Health.Status}}' my-mysql)"; \
	  sleep 2; \
	  TRIES=$$((TRIES+1)); \
	  if [ $$TRIES -ge 30 ]; then \
	    echo "Timeout waiting for MySQL healthy status"; \
	    exit 1; \
	  fi; \
	done
	@echo "Containers are healthy — verifying DB is reachable..."
	@TRIES=0; \
	until docker exec my-mysql mysqladmin ping -h localhost --silent; do \
	  echo "Waiting for MySQL ping..."; \
	  sleep 2; \
	  TRIES=$$((TRIES+1)); \
	  if [ $$TRIES -ge 30 ]; then \
	    echo "Timeout waiting for MySQL to respond to ping"; \
	    exit 1; \
	  fi; \
	done
	@echo "MySQL is reachable — starting Go app..."
	@go run cmd/main.go



migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down