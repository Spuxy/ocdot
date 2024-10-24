MODULE_DIRS = . 

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
