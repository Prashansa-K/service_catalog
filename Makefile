.PHONY: source
source:
	@. .envrc
.PHONY: run
run: source
	go run cmd/main.go

.PHONY: build
build: source
	go build -o ./.bin/service-catalog  cmd/main.go

.PHONY: fmt
fmt:
	go fmt ./...
