BIN := term-clock

all: build

.PHONY: run
run: build
	./$(BIN)

.PHONY: build
build:
	go build -o $(BIN)

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: clean
clean:
	go clean
