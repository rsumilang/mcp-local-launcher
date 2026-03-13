# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.1] - 2025-03-13

### Fixed

- golangci-lint config: use v1 schema (remove `version`, use `linters-settings` and `issues.exclude-rules`) so CI (golangci-lint 1.64.x) validates successfully

## [0.1.0] - 2025-03-13

### Added

- `open_path` — open a file or folder with the default application (e.g. folder in Finder/Explorer, file in default app)
- `reveal_in_finder` — reveal a file or folder in the system file manager and select it
- `open_with_app` — open a URL or file with a specific application (e.g. URL in Chrome, file in VS Code)
- golangci-lint configuration (`.golangci.yml`) and CI workflow that runs lint + tests on push and pull requests

## [0.0.2] - 2025-03-13

### Fixed

- Never send `id: null` in JSON-RPC error responses; use sentinel `id: -1` when request id is missing or null so strict MCP clients (e.g. Claude) accept the response

## [0.0.1] - 2025-03-13

### Added

- MCP server with `open_app` and `open_url` tools (stdio transport)
- Support for macOS, Windows, and Linux
- GitHub Actions release workflow: build binaries for linux/amd64, linux/arm64, windows/amd64, darwin/amd64, darwin/arm64 on tag push
- One-liner install script (`remote_installer.sh`) with colored output and MCP config snippets
- `.go-version` for version managers (goenv, asdf, gvm)

[Unreleased]: https://github.com/rsumilang/mcp-local-launcher/compare/v0.2.1...HEAD
[0.2.1]: https://github.com/rsumilang/mcp-local-launcher/releases/tag/v0.2.1
[0.1.0]: https://github.com/rsumilang/mcp-local-launcher/releases/tag/v0.1.0
[0.0.2]: https://github.com/rsumilang/mcp-local-launcher/releases/tag/v0.0.2
[0.0.1]: https://github.com/rsumilang/mcp-local-launcher/releases/tag/v0.0.1
