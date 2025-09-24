# Basic MCP Server in Go

A comprehensive Model Context Protocol (MCP) server implementation in Go demonstrating core MCP concepts and best practices.

## ✨ Features

### 🔧 Tools
- **`echo`**: Echoes back input text with automatic input validation
- **`calculate`**: Performs basic arithmetic operations (add, subtract, multiply, divide) with error handling
- **`system_info`**: Returns detailed system information (OS, architecture, Go version, CPU count)
- **`read_file`**: Safely reads file contents with path traversal protection and size limits

### 📚 Resources
- **`config://server`**: Server configuration in JSON format
- **`status://server`**: Real-time server status and metrics (uptime, request count, etc.)
- **`help://tools`**: Comprehensive documentation about available tools in Markdown format

### 🔒 Security Features
- Path traversal protection for file operations
- File size limits (1MB max)
- Input validation with JSON Schema
- Safe error handling

## 🚀 Quick Start

### Build and Run
```bash
go build -o mcp-server .
./mcp-server
```

### Test with Python Client
```bash
python3 test_client.py
```

## 🛠 Technical Implementation

### Architecture
- Built with the official MCP Go SDK (`github.com/modelcontextprotocol/go-sdk`)
- Uses stdio transport (standard for MCP servers)
- Modular design with separate handlers for tools and resources
- Type-safe tool arguments with automatic schema generation

### Tool Arguments
Each tool uses strongly-typed Go structs with JSON schema annotations:

```go
type EchoArgs struct {
    Text string `json:"text" jsonschema:"description:The text to echo back"`
}

type CalculateArgs struct {
    Operation string  `json:"operation" jsonschema:"description:The operation to perform,enum:add,enum:subtract,enum:multiply,enum:divide"`
    A         float64 `json:"a" jsonschema:"description:First number"`
    B         float64 `json:"b" jsonschema:"description:Second number"`
}
```

### Resource Handlers
Resources provide dynamic content with proper MIME types:

```go
s.server.AddResource(&mcp.Resource{
    URI:         "status://server",
    Name:        "Server Status",
    Description: "Current server runtime status and metrics",
    MIMEType:    "application/json",
}, s.handleStatusResource)
```

## 📖 MCP Protocol Examples

### Initialize Connection
```json
{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"capabilities": {}}}
```

### List Available Tools
```json
{"jsonrpc": "2.0", "id": 2, "method": "tools/list"}
```

### Call Echo Tool
```json
{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "echo", "arguments": {"text": "Hello, MCP!"}}}
```

### Call Calculate Tool
```json
{"jsonrpc": "2.0", "id": 4, "method": "tools/call", "params": {"name": "calculate", "arguments": {"operation": "add", "a": 15, "b": 27}}}
```

### List Resources
```json
{"jsonrpc": "2.0", "id": 5, "method": "resources/list"}
```

### Read Resource
```json
{"jsonrpc": "2.0", "id": 6, "method": "resources/read", "params": {"uri": "config://server"}}
```

## ✅ Tested Functionality

The server has been thoroughly tested and demonstrates:

- ✅ Proper MCP protocol initialization
- ✅ Tool listing and schema validation
- ✅ All 4 tools working correctly with type safety
- ✅ Resource listing and content retrieval
- ✅ Error handling and edge cases
- ✅ Request statistics and monitoring

## 📁 Project Structure

```
mcp-basic-server/
├── go.mod              # Go module definition
├── go.sum              # Dependency checksums
├── main.go             # Application entry point
├── mcp-server          # Compiled binary
├── server/
│   └── server.go       # Core MCP server implementation
├── test_client.py      # Python test client
├── sample.txt          # Sample file for testing
└── README.md           # This documentation
```

## 🔄 Next Steps

This basic implementation provides a solid foundation for understanding MCP concepts. Potential enhancements:

1. **Database Integration**: Add database tools for CRUD operations
2. **File System Tools**: Expand file operations (write, delete, list directories)
3. **API Integration**: Add tools that interact with external APIs
4. **Authentication**: Implement authentication and authorization
5. **Streaming**: Support for streaming responses
6. **Custom Transports**: Implement HTTP or WebSocket transports

## 🎯 Key Learning Outcomes

From this basic MCP server implementation, you've learned:

1. **MCP Protocol Basics**: How to implement the core MCP JSON-RPC protocol
2. **Tool Development**: Creating tools with automatic schema generation and type safety
3. **Resource Management**: Providing dynamic resources with proper content types
4. **Error Handling**: Proper error responses and validation
5. **Testing**: How to test MCP servers programmatically
6. **Go SDK Usage**: Working with the official MCP Go SDK

This foundation will help you build more complex MCP servers, including ArgoCD-specific implementations!