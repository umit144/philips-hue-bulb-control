.PHONY: build test clean

build:
	@echo "Building..."
	@go build -o bin/hue-control cmd/main.go

test:
	@echo "Running tests..."
	@go test -v ./tests/...

clean:
	@echo "Cleaning..."
	@rm -rf bin/

run:
	@echo "Running..."
	@go run cmd/main.go