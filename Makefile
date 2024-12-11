
include: .envrc

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

.PHONY: build/api
build/api:
	@echo 'Building "api" application'
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/api_linux_amd64 ./cmd/api

run/api:
	go run ./cmd/api

.PHONY: tidy
tidy:
	@echo "Tidying module dependencies"
	go mod tidy
	go mod verify

.PHONY: update
update:
	go get -u ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: security
security:
	gosec -exclude-generated ./...

.PHONY: openapi
openapi:
	go generate internal/adapter/in/http/api/server-generator.go
	go generate internal/core/domain/model/model-generator.go
	go generate internal/core/port/in/in-generator.go
