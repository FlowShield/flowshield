.PHONY: all build clean

PROG=bin/zaca
SRCS=.

# git commit hash
COMMIT_HASH=$(shell git rev-parse --short HEAD || echo "GitNotFound")
# Compilation date
BUILD_DATE=$(shell date '+%Y-%m-%d %H:%M:%S')
# Compilation conditions
CFLAGS = -ldflags "-s -w -X \"main.BuildVersion=${COMMIT_HASH}\" -X \"main.BuildDate=$(BUILD_DATE)\""

all:
	if [ ! -d "./bin/" ]; then \
	mkdir bin; \
	fi
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  $(CFLAGS) -o $(PROG) $(SRCS)

build:
	go build -race -tags=jsoniter

swagger:
	swag init

compose:
	sudo docker-compose up -d

run:
	go run main.go

clean:
	rm -rf ./bin