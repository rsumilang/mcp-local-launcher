#!/bin/sh
set -e

# ANSI colors — chosen for readability on both light and dark terminal backgrounds
R='\033[0m'           # reset
BOLD_G='\033[1;32m'  # bold green (congrats)
BOLD_C='\033[1;36m'   # bold cyan (headers)
BOLD_B='\033[1;34m'  # bold blue (path — works on light & dark)
C='\033[0;36m'       # cyan
DIM='\033[0;2m'      # dim

REPO="rsumilang/mcp-local-launcher"
BIN_DIR="${HOME}/bin"
TAG=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/;s/arm64/arm64/')
mkdir -p "$BIN_DIR"
curl -sL "https://github.com/${REPO}/releases/download/${TAG}/mcp-local-launcher-${OS}-${ARCH}" -o "${BIN_DIR}/mcp-local-launcher"
chmod +x "${BIN_DIR}/mcp-local-launcher"

# Resolve to absolute path for copy-paste
INSTALL_PATH="$(cd "$BIN_DIR" && pwd)/mcp-local-launcher"

printf '\n'
printf '%b%s%b\n' "${BOLD_C}" "==============================================" "${R}"
printf '%b  Congrats! mcp-local-launcher is installed.  %b\n' "${BOLD_G}" "${R}"
printf '%b%s%b\n' "${BOLD_C}" "==============================================" "${R}"
printf '\n'
printf '%bAdd to your MCP client using this path:%b\n' "${C}" "${R}"
printf '  %b%s%b\n' "${BOLD_B}" "$INSTALL_PATH" "${R}"
printf '\n'
printf '%b--- Claude Desktop ---%b\n' "${BOLD_C}" "${R}"
printf '%bConfig file: ~/Library/Application Support/Claude/claude_desktop_config.json (macOS)%b\n' "${DIM}" "${R}"
printf 'Add under mcpServers:\n\n'
printf '%s\n' '{' '  "local-launcher": {' "    \"command\": \"$INSTALL_PATH\"" '  }' '}'
printf '\n'
printf '%b--- Claude Code ---%b\n' "${BOLD_C}" "${R}"
printf '%bSame as above: add under mcpServers with:%b\n' "${DIM}" "${R}"
printf '"command": "%b%s%b"\n' "${BOLD_B}" "$INSTALL_PATH" "${R}"
printf '\n'
printf '%b--- Raycast ---%b\n' "${BOLD_C}" "${R}"
printf '  Name:      local-launcher\n'
printf '  Command:   %b%s%b\n' "${BOLD_B}" "$INSTALL_PATH" "${R}"
printf '  Transport: stdio\n\n'
printf '%b--- Cursor / VS Code ---%b\n' "${BOLD_C}" "${R}"
printf 'Add to MCP settings:\n\n'
printf '%s\n' '"local-launcher": {' "  \"command\": \"$INSTALL_PATH\"" '}'
printf '\n'
printf '%bEnsure ~/bin is on your PATH, or use the path above in your config.%b\n' "${DIM}" "${R}"
printf '\n'
