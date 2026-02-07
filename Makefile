build: build-web
	go build -o bin/app main.go

build-web:
	cd web && pnpm run build

format:
	gofmt -w .

lint:
	golangci-lint run

test:
	go test .

run: build-web
	go run cmd/api/main.go
