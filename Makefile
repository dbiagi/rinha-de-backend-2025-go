# Test
COVERAGE_REPORT = coverage.out
TEST_REPORT = report.out
TEST_FILES = ./internal/...
TEST_FILES_INTEGRATION = ./tests/...

# Docker
COMPOSE_FILE = ./docker/docker-compose.dev.yaml
DOCKER_COMPOSE = docker compose -f "${COMPOSE_FILE}"
COMPOSE_PROCESSOR_FILE = ./docker/docker-compose.processors.yaml
DOCKER_COMPOSE_PROCESSOR = docker compose -f "${COMPOSE_PROCESSOR_FILE}"

# Go
GOEXEC = go
GOINSTALL = ${GOEXEC} install

.PHONY: all tests test-unit test-integration test-coverage serve-dev infra-up infra-down deps build tools infra-up-processors infra-down-processors

tests:
	make test-unit
#	make test-integration

test-unit:
	@echo "Running tests..."
	${GOEXEC} test -v -coverprofile="${COVERAGE_REPORT}" ${TEST_FILES}

test-coverage:
	@echo "Resolving dependencies..."
	make deps
	@make test-unit
	# make test-integration >> ${TEST_REPORT}

test-integration:
	@echo "Running integration tests..."
	make infra-up
	${GOEXEC} test -v ${TEST_FILES_INTEGRATION}
	make infra-down

serve-dev:
	@echo "Starting server..."
	${GOEXEC} run cmd/main.go serve --env=dev

infra-up:
	@echo "Starting infrastructure..."
	${DOCKER_COMPOSE} up -d 
	make processors-infra-up
	
infra-down:
	@echo "Stopping infrastructure..."
	${DOCKER_COMPOSE} down
	make processors-infra-down

deps:
	@echo "Installing dependencies..."
	${GOEXEC} mod tidy
	${GOEXEC} mod vendor

build:
	@echo "Building..."
	${GOEXEC} build -o bin/app cmd/main.go

lint:
	@echo "Running linter..."
	golangci-lint run ./...

tools:
	${GOINSTALL} go.uber.org/mock/mockgen@latest
	${GOINSTALL} github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.5

processors-infra-up:
	${DOCKER_COMPOSE_PROCESSOR} up -d
processors-infra-down:
	${DOCKER_COMPOSE_PROCESSOR} down
