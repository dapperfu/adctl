# Makefile for adctl project
# Binary name
BIN := adctl

# Use bash as shell to ensure proper PATH handling
SHELL := /bin/bash

# Find go binary - try multiple methods
GO := $(shell command -v go 2>/dev/null || which go 2>/dev/null || echo /usr/bin/go)

# Ensure PATH includes standard locations
export PATH := /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:$(PATH)

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

# Check if go is available
GO_VERSION := $(shell $(GO) version 2>/dev/null)
ifeq ($(GO_VERSION),)
	$(error go is not available. Please install Go or ensure it's in your PATH. Tried: $(GO))
endif

.PHONY: all clean build build-test build-notest run install help

# Default target
all: build

# Clean build artifacts
clean:
	$(GO) clean -testcache
	$(GO) mod tidy
	rm -f $(BIN)
	rm -rf dist

# Build the binary (runs tests if ADCTL_HOST is set, otherwise skips)
build:
	@if [ -n "$$ADCTL_HOST" ]; then \
		echo "Running tests (ADCTL_HOST is set)..."; \
		$(GO) test ./cmd || echo "Warning: Tests failed, continuing build..."; \
	else \
		echo "Skipping tests (ADCTL_HOST not set). Use 'make build-test' to run tests."; \
	fi
	$(GO) build -o $(BIN) .

# Build with tests (requires ADCTL_HOST to be set)
build-test:
	$(GO) test ./cmd
	$(GO) build -o $(BIN) .

# Build without tests
build-notest:
	$(GO) build -o $(BIN) .

# Run the binary
run: build
	./$(BIN)

# Install the binary (builds without tests)
install: build-notest
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
	@echo "  make clean       - Clean build artifacts"
	@echo "  make build       - Build the binary (runs tests if ADCTL_HOST is set)"
	@echo "  make build-test  - Build the binary with tests (requires ADCTL_HOST)"
	@echo "  make build-notest - Build the binary without tests"
	@echo "  make run         - Build and run the binary"
	@echo "  make install     - Install the binary (builds without tests)"
	@echo "                     - With sudo: installs to /usr/local/bin"
	@echo "                     - Without sudo: installs to ~/.local/bin"
	@echo "  make help        - Display this help message"
