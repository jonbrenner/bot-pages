BINARY_NAME=bot
INSTALL_DIR=/usr/local/bin

.PHONY: build
build:
	CGO_ENABLED=0 go build -o $(BINARY_NAME)

.PHONY: install
install: build
	mkdir -p $(INSTALL_DIR)
	install -m 755 $(BINARY_NAME) $(INSTALL_DIR)
