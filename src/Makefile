fmt:
	go fmt ./...

test:
	go test ./...

generate-wire:
	wire internal/api/wire/wire.go

coverage:
	go test ./... -cover

coverage-html:
	go test ./... -cover -coverprofile coverage.out && go tool cover -html=coverage.out
