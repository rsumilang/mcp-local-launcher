#!/bin/sh
set -e
REPO="rsumilang/mcp-local-launcher"
BIN_DIR="${HOME}/bin"
TAG=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/;s/arm64/arm64/')
mkdir -p "$BIN_DIR"
curl -sL "https://github.com/${REPO}/releases/download/${TAG}/mcp-local-launcher-${OS}-${ARCH}" -o "${BIN_DIR}/mcp-local-launcher"
chmod +x "${BIN_DIR}/mcp-local-launcher"
echo "Installed to ${BIN_DIR}/mcp-local-launcher"
