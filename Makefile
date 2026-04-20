# Makefile for VPS Manager
# Usage:
#   make build   — build the binary for this machine
#   make linux   — build a Linux amd64 binary (for your VPS)
#   make clean   — remove built binaries
#   make install — install to /usr/local/bin/vps

BINARY = vps
VERSION = 0.1.0
BUILD_FLAGS = -ldflags "-X main.version=$(VERSION)"

build:
	GONOSUMDB="*" go build $(BUILD_FLAGS) -o $(BINARY) ./cmd/vps/

# Cross-compile for Linux VPS (use this if building from a Mac/Windows machine)
linux:
	GOOS=linux GOARCH=amd64 GONOSUMDB="*" go build $(BUILD_FLAGS) -o $(BINARY)-linux-amd64 ./cmd/vps/

install: build
	cp $(BINARY) /usr/local/bin/$(BINARY)
	chmod +x /usr/local/bin/$(BINARY)
	@echo "Installed to /usr/local/bin/$(BINARY)"

clean:
	rm -f $(BINARY) $(BINARY)-linux-amd64

.PHONY: build linux install clean
