all: generate bin/health
	@go mod tidy
	@go test -cover ./...
	@[ -d cmd ] && go build -ldflags "-w" -o bin/ ./cmd/...
	@go test -tags integration -c -o bin/test || true

audit:
	@which golangci-lint >/dev/null || (echo "Cannot run linters. Have you installed golangci-lint?" && false)
	@golangci-lint run

bin/health:
	@go build -ldflags "-s -w" -o bin/health gitlab.trgdev.com/gotrg/white-label/modules/health/client

generate:
	@go generate ./...

