# Variables
BINARY_NAME=botsupervisor

# Default target
all: build

# Build the Go binary
build:
	go build -o $(BINARY_NAME) main.go

kill:
	killall botsupervisor || true

# Run the Go binary (builds first)
run: build kill
	nohup ./$(BINARY_NAME) > output.log 2>&1 &

# Clean up generated files
clean:
	go clean
	rm -f $(BINARY_NAME)

# Install dependencies (optional)
deps:
	go mod tidy