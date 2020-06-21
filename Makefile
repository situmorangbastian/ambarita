BINARY=engine

# Build
engine:
	go build -o ${BINARY} cmd/ambarita/main.go

# Linter
.PHONY: lint-prepare
lint-prepare:
	@echo "Installing golangci-lint"
	@go get github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: lint
lint: vendor
	@echo "Run lint"
	@golangci-lint run ./...
