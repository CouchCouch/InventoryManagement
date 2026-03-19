build: build-web
	go build -o bin/app main.go

build-web:
	cd web && bun i && bun --bun vite build

format:
	gofumpt -w .

lint:
	golangci-lint run

test:
	go test .

run: build-web
	go run cmd/api/main.go
