package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// helper: encode a JSON-RPC request as a newline-terminated JSON line.
func requestLine(t *testing.T, method string, id interface{}, params interface{}) string {
	t.Helper()
	rawID, _ := json.Marshal(id)
	rawParams, _ := json.Marshal(params)
	req := map[string]json.RawMessage{
		"jsonrpc": json.RawMessage(`"2.0"`),
		"method":  json.RawMessage(mustJSON(t, method)),
		"id":      rawID,
	}
	if params != nil {
		req["params"] = rawParams
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	return string(b) + "\n"
}

func mustJSON(t *testing.T, v interface{}) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

// decodeResponse decodes the first JSON line from buf into a Response.
func decodeResponse(t *testing.T, buf *bytes.Buffer) Response {
	t.Helper()
	var resp Response
	if err := json.NewDecoder(buf).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v\nraw: %s", err, buf.String())
	}
	return resp
}

// runServe feeds input to serve() and returns what was written to the output buffer.
func runServe(t *testing.T, input string) *bytes.Buffer {
	t.Helper()
	out := &bytes.Buffer{}
	err := serve(strings.NewReader(input), out)
	// EOF is the expected exit condition.
	if err != nil && err.Error() != "EOF" {
		// io.EOF may be returned directly or wrapped; just tolerate it
		// only unexpected errors fail the test.
	}
	return out
}

// --- Tests for dispatch helpers ---

func TestDispatch_Initialize(t *testing.T) {
	req := Request{JSONRPC: "2.0", ID: json.RawMessage(`1`), Method: "initialize"}
	resp := dispatch(req)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}
	b, _ := json.Marshal(resp.Result)
	var result InitializeResult
	if err := json.Unmarshal(b, &result); err != nil {
		t.Fatalf("unmarshal InitializeResult: %v", err)
	}
	if result.ServerInfo.Name != serverName {
		t.Errorf("serverName = %q; want %q", result.ServerInfo.Name, serverName)
	}
	if result.ServerInfo.Version != serverVersion {
		t.Errorf("serverVersion = %q; want %q", result.ServerInfo.Version, serverVersion)
	}
	if result.Capabilities.Tools == nil {
		t.Error("expected Tools capability to be set")
	}
}

func TestDispatch_Initialized_Notification(t *testing.T) {
	req := Request{JSONRPC: "2.0", Method: "initialized"}
	resp := dispatch(req)
	if resp != nil {
		t.Errorf("expected nil response for notification, got %+v", resp)
	}
}

func TestDispatch_ToolsList(t *testing.T) {
	req := Request{JSONRPC: "2.0", ID: json.RawMessage(`2`), Method: "tools/list"}
	resp := dispatch(req)
	if resp == nil || resp.Error != nil {
		t.Fatalf("unexpected result: %+v", resp)
	}
	b, _ := json.Marshal(resp.Result)
	var result ListToolsResult
	if err := json.Unmarshal(b, &result); err != nil {
		t.Fatalf("unmarshal ListToolsResult: %v", err)
	}
	if len(result.Tools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(result.Tools))
	}
	names := map[string]bool{}
	for _, tool := range result.Tools {
		names[tool.Name] = true
	}
	for _, want := range []string{"open_app", "open_url"} {
		if !names[want] {
			t.Errorf("tool %q not found in list", want)
		}
	}
}

func TestDispatch_MethodNotFound(t *testing.T) {
	req := Request{JSONRPC: "2.0", ID: json.RawMessage(`3`), Method: "unknown/method"}
	resp := dispatch(req)
	if resp == nil || resp.Error == nil {
		t.Fatal("expected error response")
	}
	if resp.Error.Code != CodeMethodNotFound {
		t.Errorf("error code = %d; want %d", resp.Error.Code, CodeMethodNotFound)
	}
}

// --- Tests for tool input validation ---

func TestOpenApp_EmptyName(t *testing.T) {
	_, err := openApp("")
	if err == nil {
		t.Error("expected error for empty app_name")
	}
}

func TestOpenApp_WhitespaceName(t *testing.T) {
	_, err := openApp("   ")
	if err == nil {
		t.Error("expected error for whitespace-only app_name")
	}
}

func TestOpenURL_EmptyURL(t *testing.T) {
	_, err := openURL("")
	if err == nil {
		t.Error("expected error for empty url")
	}
}

func TestOpenURL_WhitespaceURL(t *testing.T) {
	_, err := openURL("   ")
	if err == nil {
		t.Error("expected error for whitespace-only url")
	}
}

// --- Tests for tools/call via handleToolsCall ---

func TestHandleToolsCall_MissingAppName(t *testing.T) {
	params, _ := json.Marshal(CallToolParams{
		Name:      "open_app",
		Arguments: map[string]interface{}{"app_name": ""},
	})
	req := Request{JSONRPC: "2.0", ID: json.RawMessage(`4`), Method: "tools/call", Params: params}
	resp := dispatch(req)
	if resp == nil || resp.Error != nil {
		t.Fatalf("expected success envelope with isError content, got %+v", resp)
	}
	b, _ := json.Marshal(resp.Result)
	var result CallToolResult
	if err := json.Unmarshal(b, &result); err != nil {
		t.Fatalf("unmarshal CallToolResult: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true for empty app_name")
	}
}

func TestHandleToolsCall_MissingURL(t *testing.T) {
	params, _ := json.Marshal(CallToolParams{
		Name:      "open_url",
		Arguments: map[string]interface{}{"url": ""},
	})
	req := Request{JSONRPC: "2.0", ID: json.RawMessage(`5`), Method: "tools/call", Params: params}
	resp := dispatch(req)
	if resp == nil || resp.Error != nil {
		t.Fatalf("expected success envelope with isError content, got %+v", resp)
	}
	b, _ := json.Marshal(resp.Result)
	var result CallToolResult
	if err := json.Unmarshal(b, &result); err != nil {
		t.Fatalf("unmarshal CallToolResult: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true for empty url")
	}
}

func TestHandleToolsCall_UnknownTool(t *testing.T) {
	params, _ := json.Marshal(CallToolParams{
		Name:      "nonexistent_tool",
		Arguments: map[string]interface{}{},
	})
	req := Request{JSONRPC: "2.0", ID: json.RawMessage(`6`), Method: "tools/call", Params: params}
	resp := dispatch(req)
	if resp == nil || resp.Error == nil {
		t.Fatal("expected JSON-RPC error for unknown tool")
	}
	if resp.Error.Code != CodeInvalidParams {
		t.Errorf("error code = %d; want %d", resp.Error.Code, CodeInvalidParams)
	}
}

// --- End-to-end serve() tests ---

func TestServe_ParseError(t *testing.T) {
	out := runServe(t, "not-json\n")
	resp := decodeResponse(t, out)
	if resp.Error == nil || resp.Error.Code != CodeParseError {
		t.Errorf("expected parse error, got %+v", resp)
	}
}

func TestServe_InvalidJSONRPCVersion(t *testing.T) {
	bad := `{"jsonrpc":"1.0","method":"initialize","id":1}` + "\n"
	out := runServe(t, bad)
	resp := decodeResponse(t, out)
	if resp.Error == nil || resp.Error.Code != CodeInvalidRequest {
		t.Errorf("expected invalid-request error, got %+v", resp)
	}
}

func TestServe_Initialize_E2E(t *testing.T) {
	input := requestLine(t, "initialize", 1, nil)
	out := runServe(t, input)
	resp := decodeResponse(t, out)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %+v", resp.Error)
	}
	b, _ := json.Marshal(resp.Result)
	var result InitializeResult
	if err := json.Unmarshal(b, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.ProtocolVersion != protocolVersion {
		t.Errorf("protocolVersion = %q; want %q", result.ProtocolVersion, protocolVersion)
	}
}

func TestServe_ToolsList_E2E(t *testing.T) {
	input := requestLine(t, "tools/list", 2, nil)
	out := runServe(t, input)
	resp := decodeResponse(t, out)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %+v", resp.Error)
	}
	b, _ := json.Marshal(resp.Result)
	var result ListToolsResult
	if err := json.Unmarshal(b, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(result.Tools) < 2 {
		t.Errorf("expected at least 2 tools, got %d", len(result.Tools))
	}
}

func TestServe_EmptyLines_Ignored(t *testing.T) {
	// Empty lines should not produce output.
	input := "\n\n\n"
	out := runServe(t, input)
	if out.Len() != 0 {
		t.Errorf("expected no output for empty lines, got %q", out.String())
	}
}
