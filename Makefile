# Makefile for adctl project
# Binary name
BIN := adctl

# Use bash as shell to ensure proper PATH handling
SHELL := /bin/bash

# Find go binary - try multiple methods
GO := $(shell command -v go 2>/dev/null || which go 2>/dev/null || echo /usr/bin/go)

# Ensure PATH includes standard locations
export PATH := /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:$(PATH)

# Install location (userspace; no sudo)
PREFIX ?= $(HOME)/.local
INSTALL_DIR ?= $(PREFIX)/bin

# Shell completion install/uninstall
# By default, we detect your login shell from $$SHELL and update the matching rc file.
COMPLETION_SHELL ?= auto
COMPLETION_MARKER := adctl completion (managed by make install)
DETECTED_LOGIN_SHELL := $(shell basename "$$SHELL" 2>/dev/null || echo bash)
BASH_COMPLETION_RC := $(HOME)/.bashrc
ZSH_COMPLETION_RC := $(HOME)/.zshrc
FISH_COMPLETION_RC := $(HOME)/.config/fish/config.fish
BASH_COMPLETION_LINE := [ -x "$(INSTALL_DIR)/$(BIN)" ] && source <("$(INSTALL_DIR)/$(BIN)" completion bash) \# $(COMPLETION_MARKER)
ZSH_COMPLETION_LINE := [ -x "$(INSTALL_DIR)/$(BIN)" ] && source <("$(INSTALL_DIR)/$(BIN)" completion zsh) \# $(COMPLETION_MARKER)
FISH_COMPLETION_LINE := test -x "$(INSTALL_DIR)/$(BIN)"; and "$(INSTALL_DIR)/$(BIN)" completion fish | source \# $(COMPLETION_MARKER)

# Check if go is available
GO_VERSION := $(shell $(GO) version 2>/dev/null)
ifeq ($(GO_VERSION),)
	$(error go is not available. Please install Go or ensure it's in your PATH. Tried: $(GO))
endif

.PHONY: all clean build build-test build-notest run install uninstall install-completion uninstall-completion help

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
	@$(MAKE) install-completion

install-completion:
	@set -eu; \
	shell="$(COMPLETION_SHELL)"; \
	if [ "$$shell" = "auto" ]; then \
		shell="$(DETECTED_LOGIN_SHELL)"; \
	fi; \
	case "$$shell" in \
		bash) rc="$(BASH_COMPLETION_RC)"; line='$(BASH_COMPLETION_LINE)';; \
		zsh) rc="$(ZSH_COMPLETION_RC)"; line='$(ZSH_COMPLETION_LINE)';; \
		fish) rc="$(FISH_COMPLETION_RC)"; line='$(FISH_COMPLETION_LINE)';; \
		*) echo "Unsupported completion shell: $$shell (set COMPLETION_SHELL=bash|zsh|fish)"; exit 1;; \
	esac; \
	mkdir -p "$$(dirname "$$rc")"; \
	touch "$$rc"; \
	if grep -Fqx "$$line" "$$rc"; then \
		echo "Completion already configured in $$rc"; \
	else \
		printf '\n%s\n' "$$line" >> "$$rc"; \
		echo "Added completion to $$rc"; \
	fi

uninstall:
	@echo "Uninstalling from $(INSTALL_DIR)"
	@rm -f "$(INSTALL_DIR)/$(BIN)"
	@$(MAKE) uninstall-completion

uninstall-completion:
	@set -eu; \
	shell="$(COMPLETION_SHELL)"; \
	if [ "$$shell" = "auto" ]; then \
		shell="$(DETECTED_LOGIN_SHELL)"; \
	fi; \
	case "$$shell" in \
		bash) rc="$(BASH_COMPLETION_RC)"; line='$(BASH_COMPLETION_LINE)';; \
		zsh) rc="$(ZSH_COMPLETION_RC)"; line='$(ZSH_COMPLETION_LINE)';; \
		fish) rc="$(FISH_COMPLETION_RC)"; line='$(FISH_COMPLETION_LINE)';; \
		*) echo "Unsupported completion shell: $$shell (set COMPLETION_SHELL=bash|zsh|fish)"; exit 1;; \
	esac; \
	if [ ! -f "$$rc" ]; then \
		echo "No rc file found at $$rc; nothing to remove"; \
		exit 0; \
	fi; \
	tmp="$$rc.tmp.$$RANDOM"; \
	grep -Fvx "$$line" "$$rc" > "$$tmp" || true; \
	mv "$$tmp" "$$rc"; \
	echo "Removed completion line from $$rc (if present)"

# Display help
help:
	@echo "Available targets:"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make build       - Build the binary (runs tests if ADCTL_HOST is set)"
	@echo "  make build-test  - Build the binary with tests (requires ADCTL_HOST)"
	@echo "  make build-notest - Build the binary without tests"
	@echo "  make run         - Build and run the binary"
	@echo "  make install     - Install the binary to $(INSTALL_DIR) and enable shell completion"
	@echo "  make uninstall   - Uninstall the binary from $(INSTALL_DIR) and remove shell completion line"
	@echo ""
	@echo "Variables:"
	@echo "  PREFIX=<dir>               - Install prefix (default: ~/.local)"
	@echo "  INSTALL_DIR=<dir>          - Install bin dir (default: $$PREFIX/bin)"
	@echo "  COMPLETION_SHELL=auto|bash|zsh|fish - Which shell rc file to update (default: auto via $$SHELL)"
	@echo "  make help        - Display this help message"
