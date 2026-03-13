# mcp-local-launcher

A minimal [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server written in Go that lets MCP clients open local desktop applications and URLs on the host machine.

Communication happens over **stdio** using JSON-RPC 2.0, making it compatible with any MCP client that supports the stdio transport (Claude Desktop, Raycast, etc.).

## Tools

| Tool | Description | Parameters |
|------|-------------|------------|
| `open_app` | Open a local application by name | `app_name` (string) |
| `open_url` | Open a URL in the default browser | `url` (string) |
| `open_path` | Open a file or folder with the default app | `path` (string) |
| `reveal_in_finder` | Reveal a file or folder in the file manager (Finder/Explorer) | `path` (string) |
| `open_with_app` | Open a URL or file with a specific application | `app_name` (string), `target` (string) |

### OS behavior

| OS | `open_app` | `open_url` | `open_path` | `reveal_in_finder` | `open_with_app` |
|----|------------|------------|------------|--------------------|-----------------|
| macOS | `open -a "<app_name>"` | `open "<url>"` | `open "<path>"` | `open -R "<path>"` | `open -a "<app>" "<target>"` |
| Windows | `powershell Start-Process "<app_name>"` | `powershell Start-Process "<url>"` | `powershell Start-Process "<path>"` | `explorer /select,<path>` | `Start-Process "<app>" -ArgumentList "<target>"` |
| Linux | exec `<app_name>` | `xdg-open "<url>"` | `xdg-open "<path>"` | `xdg-open "<parent_dir>"` | `<app_name> "<target>"` |

### Install

One-liner (macOS and Linux)—downloads the latest binary for your system into `~/bin`:

```sh
curl -sSL https://raw.githubusercontent.com/rsumilang/mcp-local-launcher/main/remote_installer.sh | sh
```

Ensure `~/bin` is on your `PATH`. Windows: download the `.exe` from the [Releases](https://github.com/rsumilang/mcp-local-launcher/releases) page.

## Requirements

- [Go 1.21+](https://go.dev/dl/)

## Releases

Pre-built binaries for **Linux** (amd64, arm64), **Windows** (amd64), and **macOS** (amd64, Apple Silicon) are published on the [Releases](https://github.com/rsumilang/mcp-local-launcher/releases) page. Download the archive for your platform, extract the binary, and (optionally) add it to your `PATH`.

To create a new release, push a version tag (e.g. `v0.1.0`). GitHub Actions will build all platforms and attach the binaries to the release:

```sh
git tag v0.1.0
git push origin v0.1.0
```

## Setup

### 1. Clone the repository

```sh
git clone https://github.com/rsumilang/mcp-local-launcher.git
cd mcp-local-launcher
```

### 2. Build

```sh
go build -o mcp-local-launcher .
```

The resulting `mcp-local-launcher` binary (or `mcp-local-launcher.exe` on Windows) is self-contained and has no external dependencies.

### 3. (Optional) Install to PATH

```sh
# macOS / Linux
sudo mv mcp-local-launcher /usr/local/bin/

# Or add the directory containing the binary to your PATH
```

## Running tests

```sh
go test ./...
```

## Lint

```sh
golangci-lint run ./...
```

Configuration is in [.golangci.yml](.golangci.yml). CI runs the linter on push and pull requests.

## MCP client configuration

### Claude Desktop

Add the server to `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) or `%APPDATA%\Claude\claude_desktop_config.json` (Windows):

```json
{
  "mcpServers": {
    "local-launcher": {
      "command": "/usr/local/bin/mcp-local-launcher"
    }
  }
}
```

Replace the `command` value with the full path to the binary if it is not on your `PATH`.

### Claude Code

Same as Claude Desktop: add a `local-launcher` entry under `mcpServers` in your Claude Code config, with `"command": "/path/to/mcp-local-launcher"`.

### Raycast

In Raycast's MCP extension settings, add a new server with:

- **Name**: `local-launcher`
- **Command**: `/usr/local/bin/mcp-local-launcher` (full path to the binary)
- **Transport**: `stdio`

### Cursor / VS Code

Add to MCP settings (Cursor: Settings → MCP; VS Code: depends on your MCP extension):

```json
"local-launcher": {
  "command": "/usr/local/bin/mcp-local-launcher"
}
```

### Generic stdio client

Any MCP client that supports the stdio transport can launch the server by running the binary directly:

```sh
mcp-local-launcher
```

The server reads JSON-RPC 2.0 messages from **stdin** (one per line) and writes responses to **stdout**. Internal log messages go to **stderr** only.

## Protocol overview

The server implements the following MCP methods:

| Method | Description |
|--------|-------------|
| `initialize` | Handshake — returns server name, version, and capabilities |
| `initialized` | Notification from client (no response sent) |
| `tools/list` | Returns the list of available tools with their JSON schemas |
| `tools/call` | Invokes a tool and returns the result |

### Example session

```jsonc
// Client → server (stdin)
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}
{"jsonrpc":"2.0","method":"initialized"}
{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"open_app","arguments":{"app_name":"Slack"}}}
{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"open_url","arguments":{"url":"https://github.com"}}}

// Server → client (stdout)
{"jsonrpc":"2.0","id":1,"result":{"protocolVersion":"2024-11-05","serverInfo":{"name":"local-launcher","version":"0.1.0"},"capabilities":{"tools":{"listChanged":false}}}}
{"jsonrpc":"2.0","id":2,"result":{"tools":[{"name":"open_app",...},{"name":"open_url",...}]}}
{"jsonrpc":"2.0","id":3,"result":{"content":[{"type":"text","text":"Opened application: Slack"}]}}
{"jsonrpc":"2.0","id":4,"result":{"content":[{"type":"text","text":"Opened URL: https://github.com"}]}}
```