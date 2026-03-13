package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	serverName      = "local-launcher"
	serverVersion   = "0.0.1"
	protocolVersion = "2024-11-05"
)

// availableTools is the static list of tools this server exposes.
var availableTools = []Tool{
	{
		Name:        "open_app",
		Description: "Open an application by name on the local machine.",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]Property{
				"app_name": {
					Type:        "string",
					Description: "Name of the application to open (e.g. \"Google Chrome\", \"Slack\").",
				},
			},
			Required: []string{"app_name"},
		},
	},
	{
		Name:        "open_url",
		Description: "Open a URL in the user's default browser.",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]Property{
				"url": {
					Type:        "string",
					Description: "The URL to open (e.g. \"https://example.com\").",
				},
			},
			Required: []string{"url"},
		},
	},
}

func main() {
	// Log internal messages to stderr only; stdout is reserved for protocol.
	log.SetOutput(os.Stderr)
	log.SetFlags(0)

	// Handle OS termination signals cleanly.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		os.Exit(0)
	}()

	if err := serve(os.Stdin, os.Stdout); err != nil && err != io.EOF {
		log.Fatalf("server error: %v", err)
	}
}

// serve reads JSON-RPC requests from r and writes responses to w until EOF.
func serve(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	// Allow very long lines (large JSON payloads).
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	encoder := json.NewEncoder(w)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			_ = encoder.Encode(errorResponse(sentinelID, CodeParseError, "parse error: "+err.Error()))
			continue
		}

		if req.JSONRPC != "2.0" {
			_ = encoder.Encode(errorResponse(req.ID, CodeInvalidRequest, "jsonrpc field must be \"2.0\""))
			continue
		}

		resp := dispatch(req)
		if resp == nil {
			// Notification – no response required.
			continue
		}
		if err := encoder.Encode(resp); err != nil {
			log.Printf("encode error: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner: %w", err)
	}
	return io.EOF
}

// dispatch routes a validated request to the appropriate handler.
func dispatch(req Request) *Response {
	switch req.Method {
	case "initialize":
		return handleInitialize(req)
	case "initialized":
		// Notification – no response.
		return nil
	case "tools/list":
		return handleToolsList(req)
	case "tools/call":
		return handleToolsCall(req)
	default:
		return errorResponse(req.ID, CodeMethodNotFound, fmt.Sprintf("method not found: %s", req.Method))
	}
}

// handleInitialize responds to the MCP "initialize" request.
func handleInitialize(req Request) *Response {
	return successResponse(req.ID, InitializeResult{
		ProtocolVersion: protocolVersion,
		ServerInfo: ServerInfo{
			Name:    serverName,
			Version: serverVersion,
		},
		Capabilities: Capabilities{
			Tools: &ToolsCapability{ListChanged: false},
		},
	})
}

// handleToolsList responds to the "tools/list" request.
func handleToolsList(req Request) *Response {
	return successResponse(req.ID, ListToolsResult{Tools: availableTools})
}

// handleToolsCall dispatches a tool invocation and returns the result.
func handleToolsCall(req Request) *Response {
	var params CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return errorResponse(req.ID, CodeInvalidParams, "invalid params: "+err.Error())
	}

	var (
		msg    string
		toolErr error
	)

	switch params.Name {
	case "open_app":
		appName, _ := params.Arguments["app_name"].(string)
		msg, toolErr = openApp(appName)
	case "open_url":
		url, _ := params.Arguments["url"].(string)
		msg, toolErr = openURL(url)
	default:
		return errorResponse(req.ID, CodeInvalidParams, fmt.Sprintf("unknown tool: %s", params.Name))
	}

	if toolErr != nil {
		return successResponse(req.ID, CallToolResult{
			Content: []ContentItem{{Type: "text", Text: toolErr.Error()}},
			IsError: true,
		})
	}
	return successResponse(req.ID, CallToolResult{
		Content: []ContentItem{{Type: "text", Text: msg}},
	})
}

// successResponse constructs a JSON-RPC 2.0 success response.
func successResponse(id json.RawMessage, result interface{}) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

// errorResponse constructs a JSON-RPC 2.0 error response.
// If id is missing or null, uses sentinelID so strict clients (e.g. Claude) accept the response.
func errorResponse(id json.RawMessage, code int, message string) *Response {
	if len(id) == 0 || string(id) == "null" {
		id = sentinelID
	}
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &RPCError{Code: code, Message: message},
	}
}
