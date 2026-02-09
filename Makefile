build:
	@go build -o bin/api cmd/api/main.go

run: build
	@./bin/api

seed:
	@go run cmd/seed/main.go

test:
	@go test -v ./...

docker:
	@echo "Building Docker image..."
	@docker build -t gotel-api .
	@echo "Running Docker container..."
	@docker run -p 3000:3000 gotel-api

clean:
	@rm -rf bin