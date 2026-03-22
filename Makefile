.PHONY: run build test test-coverage lint fmt vet check docker-build docker-run clean

## Run the application
run:
	go run cmd/server/main.go

## Build the binary
build:
	go build -o bin/server cmd/server/main.go

## Run tests
test:
	go test ./...

## Run tests with coverage
test-coverage:
	go test ./... -cover

## Run go fmt
fmt:
	go fmt ./...

## Run go vet
vet:
	go vet ./...

## Run golangci-lint
lint:
	golangci-lint run ./...

## Run all checks
check: fmt vet lint
	@echo "All checks passed!"

## Build Docker image
docker-build:
	docker build -t go-webpage-analyzer .

## Run Docker container
docker-run:
	docker run -p 8080:8080 --env-file .env go-webpage-analyzer

docker-up:
	docker-compose up --build -d

## Stop docker-compose
docker-down:
	docker-compose down

## Clean build artifacts
clean:
	rm -rf bin/