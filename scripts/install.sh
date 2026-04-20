#!/bin/bash
# VPS Manager install script
# Run this on your DigitalOcean server:
#   curl -L https://yoursite.com/install.sh | bash
#
# What it does:
#   1. Downloads the pre-built binary
#   2. Puts it in /usr/local/bin/vps
#   3. Makes it executable
#   Done. Type "vps" to launch.

set -e  # stop on any error

VERSION="0.1.0"
BINARY_URL="https://yoursite.com/releases/vps-linux-amd64"
INSTALL_PATH="/usr/local/bin/vps"

echo "Installing VPS Manager v${VERSION}..."

# Download the binary
if command -v curl &>/dev/null; then
    curl -fsSL "$BINARY_URL" -o "$INSTALL_PATH"
elif command -v wget &>/dev/null; then
    wget -q "$BINARY_URL" -O "$INSTALL_PATH"
else
    echo "Error: need curl or wget to install"
    exit 1
fi

# Make it executable
chmod +x "$INSTALL_PATH"

echo "Done! Run 'vps' to start VPS Manager."
