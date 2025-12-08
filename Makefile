.PHONY: build test clean run

BINARY_NAME=codemap
BUILD_DIR=bin

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/codemap

test:
	@echo "Running tests..."
	@go test ./...

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

run: build
	@./$(BUILD_DIR)/$(BINARY_NAME)