.PHONY: gen clean run server client docker-build

# Default build configuration
DEFAULT_OS ?= linux
DEFAULT_ARCH ?= amd64
PLATFORM ?= $(DEFAULT_OS)/$(DEFAULT_ARCH)
OS = $(word 1, $(subst /, ,$(PLATFORM)))
ARCH = $(word 2, $(subst /, ,$(PLATFORM)))

BIN_DIR = bin

gen:
	protoc --go_out=. --go-grpc_out=. proto/service.proto
	@echo "Generation complete!"

docker-build:
	docker build -t github.com/jxs1211/external:latest .

clean:
	@echo "Cleaning generated files..."
	rm -rf pb/*.go
	@echo "Clean complete!"

run:
	@echo "Running server..."
	go run server/main.go

server:
	@echo "Building server for $(OS)/$(ARCH)..."
	GOOS=$(OS) GOARCH=$(ARCH) go build -o $(BIN_DIR)/server-$(OS)-$(ARCH) ./server

client:
	@echo "Building client for $(OS)/$(ARCH)..."
	GOOS=$(OS) GOARCH=$(ARCH) go build -o $(BIN_DIR)/client-$(OS)-$(ARCH) ./client

all: gen server client
