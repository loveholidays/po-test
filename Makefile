#@ Helpers
# from https://www.thapaliya.com/en/writings/well-documented-makefiles/
help:  ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Tools
tools: ## Installs required binaries locally
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/onsi/ginkgo/v2/ginkgo

##@ Building
build-multi-arch: ## Builds po-test go binary for linux and darwin. Outputs to `bin/po-test-$GOOS-$GOARCH`.
	@echo "== build-multi-arch"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/po-test-linux-amd64
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/po-test-darwin-amd64
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o bin/po-test-darwin-arm64

build: test## Builds po-test go binary for local arch. Outputs to `bin/po-test`
	@echo "== build"
	CGO_ENABLED=0 go build -o bin/po-test

##@ Cleanup
clean: ## Deletes binaries from the bin folder
	@echo "== clean"
	rm -rfv bin/po-test-linux-amd64 bin/po-test-darwin-amd64 bin/po-test-darwin-arm64 bin/po-test

##@ Tests
test: tools ## Run unit tests
	@echo "== unit test"
	ginkgo ./...

##@ Run static checks
check: tools ## Runs lint, fmt and vet checks against the codebase
	golangci-lint run
	go fmt ./...
	go vet ./...