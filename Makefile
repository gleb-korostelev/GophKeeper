APP_NAME := gophkeeper
VERSION := 1.0.0
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

build:
	go build -ldflags "-X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE)" -o ./cmd/cli/$(APP_NAME) ./cmd/main.go

clean:
	rm -f $(APP_NAME)

.PHONY: build clean
