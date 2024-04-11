include .env
export

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

BINARY_NAME=users
CMD_PATH=./cmd/users
CONFIG_PATH=$(CMD_PATH)/config/config.yaml

MIGRATIONS_PATH=./migrations

MIGRATE_CMD=goose postgres "postgres://avito-tech-local-admin:$(DB_PASSWORD)@localhost:5433/avito-tech-local-users" up

all: clean deps build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(CMD_PATH)/main.go

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

deps:
	$(GOGET) github.com/pressly/goose/cmd/goose

migrate:
	cd $(MIGRATIONS_PATH) && $(MIGRATE_CMD)

DOCKER_COMPOSE_FILE=docker-compose.yaml

docker-up:
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d

docker-down:
	docker compose -f $(DOCKER_COMPOSE_FILE) down

docker-restart: docker-down docker-up

wait:
	sleep 5

setup: clean deps docker-up wait migrate run

run:
	$(GOBUILD) -o $(BINARY_NAME) $(CMD_PATH)/main.go && \
	./$(BINARY_NAME) -config=$(CONFIG_PATH)

test:
	$(GOTEST) -v ./...

fmt:
	$(GOCMD) fmt ./...

version:
	$(GOCMD) version

update:
	$(GOMOD) tidy

help:
	@echo "Available targets:"
	@echo "  all          - clean, deps, build"
	@echo "  build        - build the binary"
	@echo "  clean        - remove previous build"
	@echo "  deps         - get dependencies"
	@echo "  migrate      - run database migrations"
	@echo "  run          - build and run the binary"
	@echo "  test         - run tests"
	@echo "  fmt          - format source code"
	@echo "  version      - show Go version"
	@echo "  update       - update dependencies"
	@echo "  help         - show this help message"

.DEFAULT_GOAL := help
