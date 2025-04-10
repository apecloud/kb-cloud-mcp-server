.PHONY: build run test clean docker docker-run

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOTIDY=$(GOCMD) mod tidy

# Binary name
BINARY_NAME=kb-cloud-map-server
BINARY_DIR=bin

# Build the binary
build:
	mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) ./cmd/server

# Run the application
run:
	$(GORUN) ./cmd/server/main.go

# Run tests
test:
	$(GOTEST) -v ./...

# Clean up
clean:
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)

# Tidy up go.mod and go.sum
tidy:
	$(GOTIDY)

# Build the Docker image
docker:
	docker build -t apecloud/$(BINARY_NAME):latest .

# Run the Docker container
docker-run:
	docker run -p 8080:8080 \
		-e KB_CLOUD_API_KEY_NAME \
		-e KB_CLOUD_API_KEY_SECRET \
		-e KB_CLOUD_SITE \
		apecloud/$(BINARY_NAME):latest

# Run the Docker container using docker-compose
docker-compose-up:
	docker-compose up -d

# Stop the Docker container
docker-compose-down:
	docker-compose down

# Default target
all: tidy build