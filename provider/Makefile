.PHONY: start build

APP 			= zero-access
SERVER_BIN  	= ./bin/backend
GIT_COUNT 		= $(shell git rev-list --all --count)
GIT_HASH        = $(shell git rev-parse --short HEAD)
RELEASE_TAG     = $(GIT_COUNT).$(GIT_HASH)

GOPROXY=https://goproxy.oneitfarm.com,https://goproxy.cn,direct

CFLAGS = -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(SERVER_BIN)

all: release

build:
	@go build $(CFLAGS)

test:
	@go test -v $(shell go list ./...)

server:
	if [ ! -d "./bin/" ]; then \
	mkdir bin; \
	fi
	GOPROXY=$(GOPROXY) CGO_ENABLED=0 go build $(CFLAGS) ./cmd/server

release:
	if [ ! -d "./bin/" ]; then \
	mkdir bin; \
	fi
	GOPROXY=$(GOPROXY) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(CFLAGS) ./cmd/server