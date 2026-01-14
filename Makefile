# Makefile for adctl project
# Binary name
BIN := adctl

# Detect install location based on sudo usage
# If EUID is 0 or SUDO_USER is set, install to system location
# Otherwise, install to user's local bin directory
ifeq ($(shell id -u),0)
	INSTALL_DIR := /usr/local/bin
else ifdef SUDO_USER
	INSTALL_DIR := /usr/local/bin
else
	INSTALL_DIR := $(HOME)/.local/bin
endif

.PHONY: all clean build run install help

# Default target
all: build

# Clean build artifacts
clean:
	go clean -testcache
	go mod tidy
	rm -f $(BIN)
	rm -rf dist

# Build the binary
build:
	go test ./cmd
	goreleaser build --single-target --snapshot --clean
	@if [ -d dist ]; then \
		BINARY_PATH=$$(find dist -name $(BIN) -type f | head -n 1); \
		if [ -n "$$BINARY_PATH" ]; then \
			ln -fs $$BINARY_PATH ./$(BIN); \
		fi; \
	fi

# Run the binary
run: build
	./$(BIN)

# Install the binary
install: build
	@echo "Installing to $(INSTALL_DIR)"
	@mkdir -p $(INSTALL_DIR)
	@if [ -f ./$(BIN) ]; then \
		cp ./$(BIN) $(INSTALL_DIR)/$(BIN); \
		chmod +x $(INSTALL_DIR)/$(BIN); \
		echo "Installed $(BIN) to $(INSTALL_DIR)"; \
	else \
		echo "Error: Binary $(BIN) not found. Run 'make build' first."; \
		exit 1; \
	fi

# Display help
help:
	@echo "Available targets:"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make build    - Build the binary"
	@echo "  make run      - Build and run the binary"
	@echo "  make install  - Install the binary"
	@echo "                 - With sudo: installs to /usr/local/bin"
	@echo "                 - Without sudo: installs to ~/.local/bin"
	@echo "  make help     - Display this help message"
