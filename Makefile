MODULE_DIRS = . 
BINARY_NAME = "ocdot"
BINARY_VER = "v0.1"

.PHONY: clean
clean:
	rm $(BINARY_NAME)

build:
	go build .

.PHONY: lint
lint:
	@$(foreach mod,$(MODULE_DIRS), \
		(cd $(mod) && \
		echo "[lint] golangci-lint: $(mod)" && \
		golangci-lint run --path-prefix $(mod) ./...) &&) true

.PHONY: vulncheck
vulncheck:
	govulncheck .
