.PHONY: run build test clean swagger deps help

# Default target
help:
	@echo "Available commands:"
	@echo "  run      - Run the application"
	@echo "  build    - Build the application"
	@echo "  test     - Run tests"
	@echo "  swagger  - Generate Swagger documentation"
	@echo "  deps     - Download dependencies"
	@echo "  clean    - Clean build artifacts"

# Run the application
run:
	go run cmd/milkly-member/main.go

# Build the application
build:
	go build -o bin/milkly-member cmd/milkly-member/main.go

# Run tests
test:
	go test ./...

# Generate Swagger documentation
swagger:
	swag init -g cmd/milkly-member/main.go --output docs

# Download dependencies
deps:
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf docs/

# Install development tools
install-tools:
	go install github.com/swaggo/swag/cmd/swag@latest

# Docker commands (optional)
docker-build:
	docker build -t milkly-member .

docker-run:
	docker run -p 8080:8080 milkly-member