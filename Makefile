# Makefile for the aquamarine project

# Variables
BINARY_NAME=aquamarine
OUT_DEV=out/dev
OUT_PROD=out/prod

# Tools (optional). If not present, targets should still work or be skipped.
GOFUMPT?=gofumpt
GCI?=gci
GOLANGCI_LINT?=golangci-lint
GO_VET?=go vet
GO_VULNCHECK?=govulncheck

.PHONY: all build run test clean fmt lint vet check dev prod test-generated full-test

all: build

# Build the generator binary.
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) main.go
	@echo "$(BINARY_NAME) built successfully."

# Run the generator in dev mode (scaffold). Adjust as generator gains features.
run: clean
	@echo "Running generator (dev)..."
	@go run main.go generate --dev
	@echo "Done."

# Quick run in prod mode (placeholder).
prod: clean
	@echo "Running generator (prod)..."
	@go run main.go generate
	@echo "Done."

# Tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Format Go code using gofumpt and gci (fallback to go fmt if tools not available)
fmt:
	@echo "Formatting Go code..."
	@command -v $(GOFUMPT) >/dev/null 2>&1 && $(GOFUMPT) -l -w . || go fmt ./...
	@command -v $(GCI) >/dev/null 2>&1 && $(GCI) -w . || echo "gci not found; skipping import organization"
	@echo "Go code formatted."

# Run golangci-lint with comprehensive rules (best-effort if tool exists)
lint:
	@echo "Running golangci-lint..."
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 && $(GOLANGCI_LINT) run || echo "golangci-lint not found; skipping"
	@echo "Linting finished."

# Run go vet and vulnerability check
vet:
	@echo "Running go vet..."
	@$(GO_VET) ./...
	@echo "Running vulnerability check..."
	@command -v $(GO_VULNCHECK) >/dev/null 2>&1 && $(GO_VULNCHECK) ./... || echo "govulncheck not found; skipping"
	@echo "Vet checks finished."

check: fmt lint vet test
	@echo "All checks passed."

# Clean build artifacts and generated outputs
clean:
	@echo "Cleaning..."
	@rm -rf $(OUT_DEV) $(OUT_PROD) $(BINARY_NAME)
	@echo "Cleanup complete."

# Temp handy test for dual servers
test-generated:
	@echo "Testing generated dual servers..."
	@cd $(OUT_DEV) && go build
	@cd $(OUT_DEV) && timeout 5s bash -c './myapp & sleep 1 && curl localhost:8080/healthz && echo && curl localhost:8081/healthz && echo' || true

# Temp handy full regenerate + test
full-test: run test-generated
	@echo "Full test complete."
