fmt:
	go fmt ./...

test:
	go test --race ./...

clean-test-cache:
	go clean -testcache

generate-wire:
	wire internal/api/wire/wire.go

coverage:
	go test ./... -cover

coverage-html:
	go test ./... -cover -coverprofile coverage.out && go tool cover -html=coverage.out

run:
	TOKEN_EXPIRATION_TIME=30m go run ./...

build:
	go build ./...

update-packages:
	go get -u ./...
