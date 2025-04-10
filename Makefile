TARGETOS ?= $(shell go env GOOS)
TARGETARCH ?= $(shell go env GOARCH)
LDFLAGS := -s -w -extldflags "-static"

.PHONY: install-tools
install-tools: ## Install development tools
	@hash gofumpt > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go install mvdan.cc/gofumpt@latest; \
	fi

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: format
format: install-tools
	gofumpt -extra -w .

.PHONY: formatcheck
formatcheck: install-tools
	@([ -z "$(shell gofumpt -d . | head)" ]) || (echo "Source is unformatted"; exit 1)

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags '${LDFLAGS}' -v -a -tags netgo,${INSECURE_TAG} -o plugin-surge-preview .
