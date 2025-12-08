
BINARY_NAME=shellican
TEST_BINARY_NAME=shellican_test

.PHONY: all build test clean fmt lint

all: build test

build:
	go build -o $(BINARY_NAME) -v .

test:
	go test -v ./pkg/...

fmt:
	go fmt ./...

lint:
	go vet ./...
	# If you have golangci-lint installed:
	# golangci-lint run

clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(TEST_BINARY_NAME)
	rm -f *.tar.gz
