.PHONY: start build

SERVER_BIN  	= ./bin/provider
GIT_COUNT 		= $(shell git rev-list --all --count)
GIT_HASH        = $(shell git rev-parse --short HEAD)
RELEASE_TAG     = $(GIT_COUNT).$(GIT_HASH)

GOPROXY=https://goproxy.cn,direct

CFLAGS = -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(SERVER_BIN)

all: release

test:
	@go test -v $(shell go list ./...)

build:
	if [ ! -d "./bin/" ]; then \
	mkdir bin; \
	fi
	GOPROXY=$(GOPROXY) CGO_ENABLED=0 go build $(CFLAGS) ./cmd

release:
	if [ ! -d "./bin/" ]; then \
	mkdir bin; \
	fi
	GOPROXY=$(GOPROXY) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(CFLAGS) ./cmd

################################################################################
# Target: format                                                              #
################################################################################
.PHONY: format
format: modtidy
	gofumpt -l -w . && goimports -local github.com/cloudslit/ -w $(shell find ./ -type f -name '*.go' -not -path "./test")

.PHONY: modtidy
modtidy:
	go mod tidy