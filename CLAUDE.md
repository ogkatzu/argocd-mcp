# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building the MCP Server
```bash
# Build the ArgoCD MCP server (creates binary in build/)
go build -o build/argocd-mcp-server ./cmd/argocd-mcp-server

# Build with verbose output
go build -v -o build/argocd-mcp-server ./cmd/argocd-mcp-server

# Clean build cache and artifacts
go clean -cache
rm -rf build/
```

### Running the MCP Server
```bash
# Run the server directly (uses stdio transport for MCP communication)
./build/argocd-mcp-server

# Run with Go (development)
go run ./cmd/argocd-mcp-server
```

### Dependency Management
```bash
# Download and verify dependencies
go mod tidy

# View module dependencies
go mod graph

# Update dependencies to latest versions
go get -u ./...
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests in a specific package
go test -v ./internal/server

# Test ArgoCD connection (requires proper .env configuration)
go test -v ./test -run TestArgocdConnection

# Run verbose ArgoCD connection test with detailed output
go test -v ./test -run TestArgocdConnectionVerbose
```

#### ArgoCD Connection Testing
The `test/argocd_connection_test.go` file provides comprehensive testing for ArgoCD connectivity:

**Test Components:**
- **Configuration Test**: Validates environment variables and .env file loading
- **Network Connectivity**: Verifies the ArgoCD server is reachable
- **Authentication Test**: Confirms auth token validity
- **API Endpoint Test**: Tests the `/api/v1/applications` endpoint used by MCP

**Prerequisites:**
1. Configure `.env` file with valid ArgoCD credentials:
   ```bash
   ARGOCD_SERVER=https://your-argocd-server
   ARGOCD_AUTH_TOKEN=your-actual-token
   ARGOCD_INSECURE=true  # for development with self-signed certs
   ```

2. Generate ArgoCD token:
   ```bash
   argocd account generate-token --account <account-name>
   ```

**Troubleshooting:**
- **401 Errors**: Check if `ARGOCD_AUTH_TOKEN` is valid and not expired
- **Network Errors**: Verify `ARGOCD_SERVER` URL and network connectivity
- **TLS Errors**: Set `ARGOCD_INSECURE=true` for self-signed certificates

## Architecture Overview

This is an **ArgoCD MCP Server** project with dual purpose:
1. **Basic MCP Server** (`./server/` directory) - A working MCP server demonstration
2. **ArgoCD Integration** (planned) - Future ArgoCD-specific MCP server implementation

### Current Structure
```
argo_mcp/
├── cmd/argocd-mcp-server/    # Main application
│   └── main.go               # Entry point
├── internal/server/          # Private server code
│   └── server.go             # ArgoCD MCP server implementation
├── test/                     # Test files
│   └── argocd_connection_test.go # Connection tests
├── scripts/                  # Development scripts
├── docs/                     # Documentation
├── configs/                  # Configuration templates
├── build/                    # Build artifacts (gitignored)
├── go.mod                    # Go module definition
├── .mcp.json                 # MCP server configuration
└── CLAUDE.md                 # Development guide
```

### MCP Server Implementation
The MCP server (`internal/server/server.go`) is a complete working implementation that provides:

**Core Architecture Pattern:**
- **MCPServer struct**: Main server instance holding configuration, status, and MCP server
- **Handler methods**: Each tool (echo, calculate, system_info, read_file) has a dedicated handler
- **Resource providers**: Serve configuration, status, and help documentation
- **Transport**: Uses stdio transport for communication with MCP clients

**Key Go Patterns Used:**
- Method receivers for server functionality: `func (s *MCPServer) methodName()`
- Pointer types for efficiency: `*MCPServer`, `*ServerConfig`
- Multiple return values for error handling: `(result, any, error)`
- JSON struct tags for serialization: `json:"field_name"`
- Context-based request handling: `context.Context` parameters

### MCP Integration with Claude Code
- **Configuration**: `.mcp.json` defines the server as "argocd-mcp"
- **Binary Path**: Points to `./build/argocd-mcp-server` (built from `go build`)
- **Transport**: Uses stdio (standard input/output) for MCP protocol communication
- **Resources Available**: ArgoCD applications list

## Go Development Context

### Module Structure
- **Main module**: `argo_mcp` (root level, currently minimal)
- **Server module**: Self-contained in root directory
- **Go version**: 1.25.1
- **Key dependency**: `github.com/modelcontextprotocol/go-sdk v0.5.0`

### Code Conventions
- Use Context7 MCP for Go language reference and best practices
- Follow standard Go conventions: CamelCase for exported types, camelCase for unexported
- Error handling: Return errors as values, check explicitly
- Struct initialization: Use `&StructName{field: value}` pattern
- JSON tags: Use snake_case for JSON field names

### Environment Variables (Future ArgoCD Integration)
- `ARGOCD_SERVER`: ArgoCD server URL
- `ARGOCD_AUTH_TOKEN`: Authentication token
- `ARGOCD_INSECURE`: Skip TLS verification for development
- `MCP_TRANSPORT`: Transport method (stdio, http, websocket)

## Project Status

**Current State**: The basic MCP server is fully functional and integrated with Claude Code
**Future Plans**: ArgoCD integration features are planned but not yet implemented

The current codebase serves as both a learning project for Go development and a foundation for future ArgoCD integration.