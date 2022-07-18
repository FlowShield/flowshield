NAME := client
GIT_COUNT 		= $(shell git rev-list --all --count)
GIT_HASH        = $(shell git rev-parse --short HEAD)
RELEASE_TAG     = $(RELEASE_VERSION).$(GIT_COUNT).$(GIT_HASH)

BUILD_DIR := bin
CMD_DIR := ./cmd
BUILD = env CGO_ENABLED=0 $(GO_DIR)go build -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(BUILD_DIR)
GOBUILD = $(BUILD)/$@/$(NAME) $(CMD_DIR)

.PHONY: client release test
normal: clean client

clean:
	rm -rf $(BUILD_DIR)
	rm -f machine.lock

client:
	mkdir -p $(BUILD_DIR)
	$(BUILD)/$(NAME) $(CMD_DIR)

darwin-amd64:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=darwin $(GOBUILD)

darwin-arm64:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=arm64 GOOS=darwin $(GOBUILD)

linux-386:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=386 GOOS=linux $(GOBUILD)

linux-amd64:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=linux $(GOBUILD)

linux-arm:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=arm GOOS=linux $(GOBUILD)

linux-armv5:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=arm GOOS=linux GOARM=5 $(GOBUILD)

linux-armv6:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=arm GOOS=linux GOARM=6 $(GOBUILD)

linux-armv7:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=arm GOOS=linux GOARM=7 $(GOBUILD)

linux-armv8:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=arm64 GOOS=linux $(GOBUILD)

linux-mips-softfloat:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=mips GOMIPS=softfloat GOOS=linux $(GOBUILD)

linux-mips-hardfloat:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=mips GOMIPS=hardfloat GOOS=linux $(GOBUILD)

linux-mipsle-softfloat:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=mipsle GOMIPS=softfloat GOOS=linux $(GOBUILD)

linux-mipsle-hardfloat:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=mipsle GOMIPS=hardfloat GOOS=linux $(GOBUILD)

linux-mips64:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=mips64 GOOS=linux $(GOBUILD)

linux-mips64le:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=mips64le GOOS=linux $(GOBUILD)

freebsd-386:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=386 GOOS=freebsd $(GOBUILD)

freebsd-amd64:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=freebsd $(GOBUILD)

windows-386:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=386 GOOS=windows $(GOBUILD)

windows-amd64:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=windows $(GOBUILD)

windows-arm:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=arm GOOS=windows $(GOBUILD)

windows-armv6:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=arm GOOS=windows GOARM=6 $(GOBUILD)

windows-armv7:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=arm GOOS=windows GOARM=7 $(GOBUILD)

windows-arm64:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=arm64 GOOS=windows $(GOBUILD)
