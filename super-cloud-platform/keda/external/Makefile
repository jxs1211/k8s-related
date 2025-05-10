.PHONY: gen clean run server client

gen:
	protoc --go_out=. --go-grpc_out=. proto/service.proto
	@echo "Generation complete!"

clean:
	@echo "Cleaning generated files..."
	rm -rf pb/*.go
	@echo "Clean complete!"

run:
	@echo "Running server..."
	go run server/main.go

server:
	@echo "Building server..."
	go build -o bin/server ./server

client:
	@echo "Building client..."
	go build -o bin/client ./client

all: gen server client