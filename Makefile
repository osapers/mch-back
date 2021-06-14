build:
	@go build -ldflags "-w -s" -o back

gomod:
	@go mod tidy
	@go mod vendor

lint:
	@golangci-lint run

code-generate:
	@go generate ./...

pre-commit: gomod code-generate lint