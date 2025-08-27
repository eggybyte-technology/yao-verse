.PHONY: all clean generate lint breaking format deps

# Protocol buffer generation and management
all: deps generate

# Install required tools
deps:
	@echo "Installing buf and dependencies..."
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate Go code from proto files
generate:
	@echo "Generating Go code from proto files..."
	@buf generate
	@echo "✓ Proto generation completed"

# Lint proto files
lint:
	@echo "Linting proto files..."
	@buf lint
	@echo "✓ Proto linting completed"

# Check for breaking changes
breaking:
	@echo "Checking for breaking changes..."
	@buf breaking --against '.git#branch=main'
	@echo "✓ Breaking change check completed"

# Format proto files
format:
	@echo "Formatting proto files..."
	@buf format -w
	@echo "✓ Proto formatting completed"

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	@rm -rf ../shared/gen/proto
	@rm -rf ../docs/api-reference/*.swagger.json
	@echo "✓ Cleanup completed"

# Update dependencies
update:
	@echo "Updating buf dependencies..."
	@buf mod update
	@echo "✓ Dependencies updated"

# Validate proto files
validate: lint
	@echo "Validating proto files..."
	@buf build
	@echo "✓ Proto validation completed"

# Development workflow: format, lint, generate
dev: format lint generate

# CI workflow: lint, breaking, generate  
ci: lint breaking generate

# Help
help:
	@echo "Available targets:"
	@echo "  all        - Install dependencies and generate code"
	@echo "  deps       - Install required tools"
	@echo "  generate   - Generate Go code from proto files"
	@echo "  lint       - Lint proto files"
	@echo "  breaking   - Check for breaking changes"
	@echo "  format     - Format proto files"
	@echo "  clean      - Clean generated files"
	@echo "  update     - Update buf dependencies"
	@echo "  validate   - Validate proto files"
	@echo "  dev        - Development workflow (format, lint, generate)"
	@echo "  ci         - CI workflow (lint, breaking, generate)"
	@echo "  help       - Show this help message" 