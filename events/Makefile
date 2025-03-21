BINARY_NAME = auth
BUILD_DIR = bin
INTEGRATION_TAGS = integration

.PHONY: build run test test-integration clean

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd

run: build
	@echo "Starting server..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

test:
	@echo "Running unit tests..."
	@go test -v -cover ./...

test-integration:
	@echo "Running integration tests..."
	@go test -v -tags=$(INTEGRATION_TAGS) ./...

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)