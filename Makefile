.PHONY: build test vet lint clean check-deps

build:
	go build ./...

test:
	go test -race ./...

vet:
	go vet ./...

lint:
	staticcheck ./...

# enforces the dependency rule: domain must never import adapters
check-deps:
	@echo "Checking dependency rule..."
	@if go list -f '{{.ImportPath}}: {{.Imports}}' ./internal/domain/... | grep -q 'internal/adapters'; then \
		echo "FAIL: internal/domain imports internal/adapters"; exit 1; \
	fi
	@echo "OK"

clean:
	rm -rf bin/ dist/ coverage.txt

ci: vet lint test check-deps
