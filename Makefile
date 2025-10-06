BINARY_NAME=machina
CMD_DIR=./cmd/machina

build:
	@echo "Building $(BINARY_NAME) for host..."
	go build -o bin/$(BINARY_NAME) $(CMD_DIR)

build-windows:
	@echo "Cross-compiling for Windows amd64..."
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME).exe $(CMD_DIR)

build-all: clean build build-windows

clean:
	@echo "Cleaning..."
	rm -rf bin