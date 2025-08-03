.PHONY: build run clean docker-build docker-run

# Build the application
build:
	go build -o badger-web-ui .

# Run the application
run:
	go run .

# Clean build artifacts
clean:
	rm -f badger-web-ui
	rm -rf badger-data

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run tests
test:
	go test -v ./...

# Build Docker image
docker-build:
	docker build -t badger-web-ui .

# Run Docker container
docker-run:
	docker run -p 8080:8080 -v $(PWD)/badger-data:/root/badger-data badger-web-ui

# Development mode with auto-reload (requires air)
dev: tools
	air

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint: tools
	golangci-lint run

# Download development tools (golangci-lint and air)
tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/air-verse/air@latest