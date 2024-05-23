.PHONY: source
source:
	@. .envrc
.PHONY: run
run: source init-localdev
	go run cmd/main.go

.PHONY: init-localdev
init-localdev:
	./scripts/init-localdev.sh

.PHONY: build
build: source
	go build -o ./.bin/service-catalog  cmd/main.go

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test ./...
