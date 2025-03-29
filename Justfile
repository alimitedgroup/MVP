set export := true

# List all recipes
list:
    @just -l

# Build the docker images
build:
    docker compose pull
    docker compose build

# Start the docker-compose setup
up:
    docker compose up -d --build

# Stop the docker-compose setup
down:
    docker compose down

# Reset the containers
reset:
    docker compose down -v --remove-orphans
    just up

# Run `go generate`
generate:
    go generate ./...

# Run all tests
test: generate
    go test ./...

coverage: generate
    go test -race -covermode atomic ./...

# Run go vet
vet: generate
    go vet ./...

# Run `staticcheck`
staticcheck: generate
    go run honnef.co/go/tools/cmd/staticcheck@latest ./...

# Run `golangci-lint`
golangci-lint: generate
    go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5 run

# Run `gocyclo`
gocyclo: generate
    #!/bin/bash
    if [ $(go run github.com/fzipp/gocyclo/cmd/gocyclo@latest -over 15 . | wc -l) -gt 0 ]; then
        echo "Complexity too high in some functions"
        go run github.com/fzipp/gocyclo/cmd/gocyclo@latest -over 15 .
        exit 1
    fi

# Run all checks that are run in CI
ci: test vet staticcheck golangci-lint gocyclo
