GO ?= go
output ?= dbg-go
TEST_FLAGS ?= -v -race -tags=test_all
LDFLAGS := -X github.com/falcosecurity/driverkit/pkg/driverbuilder/builder.defaultImageTag=v0.15.0

.PHONY: build
build: ${output}

.PHONY: ${output}
${output}:
	CGO_ENABLED=0 $(GO) build -ldflags '${LDFLAGS}' -o $@

.PHONY: clean
clean:
	$(RM) -R ${output}

.PHONY: test
test:
	$(GO) test ${TEST_FLAGS} ./...
