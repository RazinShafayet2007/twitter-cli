#!/bin/bash
# Twitter CLI Uninstaller

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

INSTALL_DIR="$HOME/.twitter-cli"

echo -e "${YELLOW}Twitter CLI Uninstaller${NC}"
echo "======================="
echo ""

if [ ! -d "$INSTALL_DIR" ]; then
    echo "Twitter CLI is not installed."
    exit 0
fi

# Confirm
read -p "This will remove Twitter CLI and all data. Continue? (y/N) " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled."
    exit 0
fi

echo "Removing installation directory..."
rm -rf "$INSTALL_DIR"

echo -e "${GREEN}✓${NC} Removed $INSTALL_DIR"
echo ""

# Remove from PATH
SHELL_NAME=$(basename "$SHELL")
case $SHELL_NAME in
    bash) RC_FILE="$HOME/.bashrc" ;;
    zsh) RC_FILE="$HOME/.zshrc" ;;
    fish) RC_FILE="$HOME/.config/fish/config.fish" ;;
    *) RC_FILE="$HOME/.profile" ;;
esac

if [ -f "$RC_FILE" ] && grep -q "twitter-cli/bin" "$RC_FILE"; then
    echo "Removing from $RC_FILE..."
    # Create backup
    cp "$RC_FILE" "${RC_FILE}.backup"
    # Remove lines
    sed -i '/# Twitter CLI/d' "$RC_FILE"
    sed -i '/twitter-cli\/bin/d' "$RC_FILE"
    echo -e "${GREEN}✓${NC} Removed from PATH"
    echo ""
    echo "Backup saved to: ${RC_FILE}.backup"
fi

echo ""
echo -e "${GREEN}Uninstall complete!${NC}"
echo ""
echo "Reload your shell:"
echo "  source $RC_FILE"