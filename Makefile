GO ?= go
output ?= dbg-go
TEST_FLAGS ?= -v -race -tags=test_all

.PHONY: build
build: ${output}

.PHONY: ${output}
${output}:
	$(GO) build -o $@

.PHONY: clean
clean:
	$(RM) -R ${output}

.PHONY: test
test:
	$(GO) test ${TEST_FLAGS} ./...
