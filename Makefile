.PHONY: run
run: init-localdev
	go run cmd/main.go

.PHONY: init-localdev
init-localdev:
	./scripts/init-localdev.sh

.PHONY: delete-localdev
delete-localdev:
	./scripts/delete-localdev.sh

.PHONY: build
build:
	go build -o ./.bin/service-catalog  cmd/main.go

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test ./...
