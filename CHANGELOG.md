# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/rsumilang/mcp-local-launcher/compare/v0.0.2...HEAD
[0.0.2]: https://github.com/rsumilang/mcp-local-launcher/releases/tag/v0.0.2
[0.0.1]: https://github.com/rsumilang/mcp-local-launcher/releases/tag/v0.0.1
