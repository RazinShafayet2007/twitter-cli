#!/bin/bash
# Twitter CLI Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/YOUR_USERNAME/twitter-cli/main/scripts/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
INSTALL_DIR="$HOME/.twitter-cli/bin"
CONFIG_DIR="$HOME/.twitter-cli"
BINARY_NAME="twt"
REPO="YOUR_USERNAME/twitter-cli"
VERSION="${VERSION:-latest}"

echo -e "${GREEN}Twitter CLI Installer${NC}"
echo "====================="
echo ""

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
    *)
        echo -e "${RED}Error: Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

case $OS in
    linux|darwin) ;;
    mingw*|msys*|cygwin*) OS="windows" ;;
    *)
        echo -e "${RED}Error: Unsupported OS: $OS${NC}"
        exit 1
        ;;
esac

echo "Detected: ${OS}/${ARCH}"
echo ""

# Create installation directory
echo "Creating installation directory..."
mkdir -p "$INSTALL_DIR"
mkdir -p "$CONFIG_DIR"

# Determine download URL
if [ "$VERSION" = "latest" ]; then
    DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/twt-${OS}-${ARCH}"
else
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/twt-${OS}-${ARCH}"
fi

if [ "$OS" = "windows" ]; then
    DOWNLOAD_URL="${DOWNLOAD_URL}.exe"
    BINARY_NAME="${BINARY_NAME}.exe"
fi

echo "Downloading from: $DOWNLOAD_URL"
echo ""

# Download binary
if command -v curl &> /dev/null; then
    curl -fsSL "$DOWNLOAD_URL" -o "$INSTALL_DIR/$BINARY_NAME"
elif command -v wget &> /dev/null; then
    wget -q "$DOWNLOAD_URL" -O "$INSTALL_DIR/$BINARY_NAME"
else
    echo -e "${RED}Error: curl or wget is required${NC}"
    exit 1
fi

# Make executable
chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo -e "${GREEN}✓${NC} Binary installed to $INSTALL_DIR/$BINARY_NAME"
echo ""

# Check if already in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "${YELLOW}Adding to PATH...${NC}"
    
    # Detect shell
    SHELL_NAME=$(basename "$SHELL")
    
    case $SHELL_NAME in
        bash)
            RC_FILE="$HOME/.bashrc"
            ;;
        zsh)
            RC_FILE="$HOME/.zshrc"
            ;;
        fish)
            RC_FILE="$HOME/.config/fish/config.fish"
            ;;
        *)
            RC_FILE="$HOME/.profile"
            ;;
    esac
    
    # Add to PATH
    if [ -f "$RC_FILE" ]; then
        if ! grep -q "twitter-cli/bin" "$RC_FILE"; then
            echo "" >> "$RC_FILE"
            echo "# Twitter CLI" >> "$RC_FILE"
            echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$RC_FILE"
            echo -e "${GREEN}✓${NC} Added to $RC_FILE"
            echo ""
            echo -e "${YELLOW}Run this to update your current shell:${NC}"
            echo "  source $RC_FILE"
            echo ""
            echo "Or close and reopen your terminal."
        fi
    else
        echo -e "${YELLOW}Could not find shell config file.${NC}"
        echo "Add this to your shell's RC file:"
        echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
    fi
else
    echo -e "${GREEN}✓${NC} Already in PATH"
fi

echo ""
echo -e "${GREEN}Installation complete!${NC}"
echo ""
echo "Try it out:"
echo "  $BINARY_NAME --help"
echo "  $BINARY_NAME user create alice"
echo ""
echo "Documentation: https://github.com/${REPO}"