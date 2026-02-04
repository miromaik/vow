#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO="miromaik/vow"
BINARY_NAME="vow"
INSTALL_DIR="/usr/local/bin"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo -e "${RED}Unsupported architecture: $ARCH${NC}" && exit 1 ;;
esac

echo -e "${GREEN}VOW Installer${NC}"
echo "================================"
echo "OS: $OS"
echo "Architecture: $ARCH"
echo ""

# Check if Go is installed
if command -v go &> /dev/null; then
    echo -e "${YELLOW}Installing via 'go install'...${NC}"
    go install github.com/${REPO}@latest
    
    # Check if $GOPATH/bin is in PATH
    if [[ ":$PATH:" != *":$GOPATH/bin:"* ]] && [[ ":$PATH:" != *":$HOME/go/bin:"* ]]; then
        echo -e "${YELLOW}Warning: Add Go bin directory to your PATH:${NC}"
        echo "  export PATH=\$PATH:\$(go env GOPATH)/bin"
        echo ""
        echo "Add this line to your ~/.zshrc or ~/.bashrc"
    fi
    
    echo -e "${GREEN}✓ Installed successfully via go install${NC}"
    echo "Run: vow"
else
    echo -e "${YELLOW}Go not found. Building from source...${NC}"
    
    # Check if git is installed
    if ! command -v git &> /dev/null; then
        echo -e "${RED}Error: git is required${NC}"
        exit 1
    fi
    
    # Create temp directory
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"
    
    echo "Cloning repository..."
    git clone --depth 1 https://github.com/${REPO}.git
    cd vow
    
    echo "Building binary..."
    go build -o "$BINARY_NAME" .
    
    # Install binary
    echo "Installing to $INSTALL_DIR..."
    if [ -w "$INSTALL_DIR" ]; then
        mv "$BINARY_NAME" "$INSTALL_DIR/"
    else
        echo -e "${YELLOW}Requesting sudo access to install to $INSTALL_DIR${NC}"
        sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
    fi
    
    # Cleanup
    cd
    rm -rf "$TMP_DIR"
    
    echo -e "${GREEN}✓ Installed successfully to $INSTALL_DIR/$BINARY_NAME${NC}"
    echo "Run: vow"
fi

echo ""
echo "================================"
echo -e "${GREEN}Installation complete!${NC}"
