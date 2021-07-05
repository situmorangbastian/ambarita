BINARY=engine

# Build
engine:
	go build -o ${BINARY} cmd/ambarita/main.go

# Dependency Management
.PHONY: vendor
vendor: go.mod go.sum
	@GO111MODULE=on go get ./...

# Linter
.PHONY: lint-prepare
lint-prepare: vendor
	@echo "Installing golangci-lint"
	@go get github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: lint
lint: vendor
	@echo "Run lint"
	@golangci-lint run ./...
