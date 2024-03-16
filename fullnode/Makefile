GOPROXY=https://goproxy.cn,direct

PROG=bin/fullnode


SRCS=.

# installation
INSTALL_PREFIX=/usr/local/fullnode

# Configure the installation directory
CONF_INSTALL_PREFIX=/usr/local/fullnode

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
	GOPROXY=$(GOPROXY) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  $(CFLAGS) -o $(PROG) $(SRCS)

# Compile race version
race:
	if [ ! -d "./bin/" ]; then \
    	mkdir bin; \
    	fi
	go build $(CFLAGS) -race -o $(PROG) $(SRCS)

# release version
RELEASE_DATE = $(shell date '+%Y%m%d%H%M%S')
RELEASE_VERSION = $(shell git rev-parse --short HEAD || echo "GitNotFound")
RELEASE_DIR=release_bin
RELEASE_BIN_NAME=fullnode
release:
	if [ ! -d "./$(RELEASE_DIR)/$(RELEASE_DATE)_$(RELEASE_VERSION)" ]; then \
	mkdir ./$(RELEASE_DIR)/$(RELEASE_DATE)_$(RELEASE_VERSION); \
	fi
	go build  $(CFLAGS) -o $(RELEASE_DIR)/$(RELEASE_DATE)_$(RELEASE_VERSION)/$(RELEASE_BIN_NAME)_linux_amd64 $(SRCS)

install:
	cp $(PROG) $(INSTALL_PREFIX)/bin

	if [ ! -d "${CONF_INSTALL_PREFIX}" ]; then \
	mkdir $(CONF_INSTALL_PREFIX); \
	fi

	cp -R config/* $(CONF_INSTALL_PREFIX)

clean:
	rm -rf ./bin

	rm -rf $(INSTALL_PREFIX)/bin/fullnode

	rm -rf $(CONF_INSTALL_PREFIX)

run:
	go run main.go

run_race:
	go run --race main.go
