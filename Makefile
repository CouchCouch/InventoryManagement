build: build-web
	go build -o bin/app cmd/api/main.go

build-web:
	cd web && bun i && bun --bun vite build

fmt:
	gofumpt -w .

lint:
	golangci-lint run

test:
	go test ./...

testv:
	go test -v ./...

run: build-web
	go run cmd/api/main.go
