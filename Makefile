.PHONY: build install uninstall clean test

BINARY_NAME=vow
INSTALL_DIR=/usr/local/bin

build:
	go build -o $(BINARY_NAME) .

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@if [ -w $(INSTALL_DIR) ]; then \
		mv $(BINARY_NAME) $(INSTALL_DIR)/; \
	else \
		sudo mv $(BINARY_NAME) $(INSTALL_DIR)/; \
	fi
	@echo "✓ Installed successfully"

uninstall:
	@echo "Removing $(BINARY_NAME) from $(INSTALL_DIR)..."
	@if [ -w $(INSTALL_DIR)/$(BINARY_NAME) ]; then \
		rm $(INSTALL_DIR)/$(BINARY_NAME); \
	else \
		sudo rm $(INSTALL_DIR)/$(BINARY_NAME); \
	fi
	@echo "✓ Uninstalled successfully"

clean:
	@rm -f $(BINARY_NAME)
	@rm -f REPORT.txt
	@echo "✓ Cleaned build artifacts"

test:
	go test ./...

cross-compile:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 .
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe .
	@echo "✓ Cross-compilation complete"
