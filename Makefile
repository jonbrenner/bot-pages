BINARY_NAME=bot
INSTALL_DIR=/usr/local/bin
ARCH=$(shell uname -m)
VERSION=$(shell git describe --tags --always)

.PHONY: build
build:
	CGO_ENABLED=0 go build -o $(BINARY_NAME)

.PHONY: install
install: build
	mkdir -p $(INSTALL_DIR)
	install -m 755 $(BINARY_NAME) $(INSTALL_DIR)

.PHONY: release
release: build
	mkdir -p release
	tar -czf release/$(BINARY_NAME)-$(VERSION)-$(ARCH).tar.gz --no-same-owner $(BINARY_NAME)