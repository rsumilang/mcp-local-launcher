# mcp-local-launcher

A minimal [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server written in Go that lets MCP clients open local desktop applications and URLs on the host machine.

Communication happens over **stdio** using JSON-RPC 2.0, making it compatible with any MCP client that supports the stdio transport (Claude Desktop, Raycast, etc.).

## Tools

| Tool | Description | Parameter |
|------|-------------|-----------|
| `open_app` | Open a local application by name | `app_name` (string) |
| `open_url` | Open a URL in the default browser | `url` (string) |

### OS behavior

| OS | `open_app` | `open_url` |
|----|-----------|-----------|
| macOS | `open -a "<app_name>"` | `open "<url>"` |
| Windows | `powershell Start-Process "<app_name>"` | `powershell Start-Process "<url>"` |
| Linux | exec `<app_name>` directly | `xdg-open "<url>"` |

## Requirements

- [Go 1.21+](https://go.dev/dl/)

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

### Raycast

In Raycast's MCP extension settings, add a new server with:

- **Name**: `local-launcher`
- **Command**: `/usr/local/bin/mcp-local-launcher` (full path to the binary)
- **Transport**: `stdio`

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