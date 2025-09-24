# ArgoCD MCP Server

A Model Context Protocol (MCP) server for ArgoCD integration, enabling seamless GitOps operations through Claude Code.

## ✨ Features

### 🚀 ArgoCD Integration
- **Applications Resource**: List all ArgoCD applications with status information
- **Connection Testing**: Built-in connectivity and authentication testing
- **Secure Communication**: TLS configuration with optional insecure mode for development

### 🔧 Architecture
- **Standard Go Layout**: Following Go project conventions with `cmd/`, `internal/`, `test/` structure
- **MCP Protocol**: Full Model Context Protocol implementation using the official Go SDK
- **Environment Configuration**: Flexible configuration via environment variables and `.env` files

## 🚀 Quick Start

### Prerequisites
- Go 1.25.1 or later
- Access to an ArgoCD server
- ArgoCD authentication token

### 1. Clone and Setup
```bash
git clone <repository-url>
cd argo_mcp
cp configs/.env.example .env
```

### 2. Configure ArgoCD Connection
Edit `.env` with your ArgoCD details:
```bash
ARGOCD_SERVER=https://your-argocd-server
ARGOCD_AUTH_TOKEN=your-auth-token
ARGOCD_INSECURE=true  # for development with self-signed certs
```

### 3. Generate ArgoCD Token
```bash
argocd account generate-token --account <account-name>
```

### 4. Build and Test
```bash
# Build the server
go build -o build/argocd-mcp-server ./cmd/argocd-mcp-server

# Test connection
go test -v ./test -run TestArgocdConnectionVerbose

# Run the server
./build/argocd-mcp-server
```

## 🧪 Testing

### Connection Testing
Test your ArgoCD connection before using the MCP server:

```bash
# Basic connection test
go test -v ./test -run TestArgocdConnection

# Verbose diagnostic test
go test -v ./test -run TestArgocdConnectionVerbose
```

The verbose test provides detailed output:
- ✅ Configuration validation
- ✅ Network connectivity check
- ✅ Authentication verification
- ✅ API endpoint testing

### Integration Testing
```bash
# Run all tests
go test ./...

# Test specific components
go test -v ./internal/server
```

## 📁 Project Structure

```
argo_mcp/
├── cmd/argocd-mcp-server/    # Main application
│   └── main.go               # Entry point
├── internal/server/          # Private server code
│   └── server.go             # ArgoCD MCP server implementation
├── test/                     # Test files
│   ├── argocd_connection_test.go # Connection tests
│   └── integration/          # Integration tests
├── scripts/                  # Development scripts
│   ├── test_client.py        # Python test client
│   └── verify_config.py      # Configuration verification
├── docs/                     # Documentation
│   └── README.md             # Additional documentation
├── configs/                  # Configuration templates
│   ├── .env.example          # Environment template
│   └── kind-config.yaml      # Kubernetes config
├── build/                    # Build artifacts
├── .env                      # Your environment config
├── .mcp.json                 # MCP server configuration
├── go.mod                    # Go module definition
└── CLAUDE.md                 # Development guide
```

## 🔧 Development Commands

### Building
```bash
# Build the ArgoCD MCP server
go build -o build/argocd-mcp-server ./cmd/argocd-mcp-server

# Build with verbose output
go build -v -o build/argocd-mcp-server ./cmd/argocd-mcp-server
```

### Running
```bash
# Run the server directly
./build/argocd-mcp-server

# Run with Go (development)
go run ./cmd/argocd-mcp-server
```

### Dependencies
```bash
# Download and verify dependencies
go mod tidy

# Update dependencies to latest versions
go get -u ./...
```

## 🔌 Claude Code Integration

### MCP Configuration
The server integrates with Claude Code via `.mcp.json`:

```json
{
  "mcpServers": {
    "argocd-mcp": {
      "command": "/Users/skatz/argo_mcp/build/argocd-mcp-server",
      "args": []
    }
  }
}
```

### Available Resources
- **`argocd://applications`**: List all ArgoCD applications with metadata, status, and health information

## 🛠 Technical Details

### MCP Server Implementation
- Built with the official MCP Go SDK (`github.com/modelcontextprotocol/go-sdk`)
- Uses stdio transport for Claude Code communication
- Implements ArgoCD REST API client with configurable TLS settings
- Environment-based configuration with `.env` support

### ArgoCD API Integration
```go
type ArgocdApplication struct {
    Metadata struct {
        Name      string `json:"name"`
        Namespace string `json:"namespace"`
    } `json:"metadata"`
    Status struct {
        Sync struct {
            Status string `json:"status"`
        } `json:"sync"`
        Health struct {
            Status string `json:"status"`
        } `json:"health"`
    } `json:"status"`
}
```

## 🔒 Security

- **Authentication**: Bearer token authentication with ArgoCD
- **TLS Configuration**: Configurable TLS verification (disable for development)
- **Error Handling**: Secure error responses without credential leakage
- **Environment Isolation**: Credentials stored in environment variables

## 📖 Usage Examples

### List ArgoCD Applications
Once connected to Claude Code with `/mcp`, you can:

1. List all applications: Access the `argocd://applications` resource
2. View application status, health, and sync information
3. Get detailed metadata for each application

### Troubleshooting

**401 Authentication Errors:**
- Verify `ARGOCD_AUTH_TOKEN` is valid and not expired
- Check token has appropriate permissions
- Run connection test: `go test -v ./test -run TestArgocdConnectionVerbose`

**Network Connectivity Issues:**
- Verify `ARGOCD_SERVER` URL is correct and reachable
- Check firewall settings
- For development: Set `ARGOCD_INSECURE=true` for self-signed certificates

**MCP Connection Issues:**
- Rebuild the server after configuration changes
- Restart Claude Code to reload MCP server
- Check `.mcp.json` points to correct binary path

## 🚧 Future Enhancements

- **Application Management**: Create, update, delete applications
- **Sync Operations**: Trigger application synchronization
- **Rollback Operations**: Rollback to previous versions
- **Cluster Management**: Multi-cluster support
- **Real-time Updates**: WebSocket support for live status updates
- **CLI Integration**: ArgoCD CLI command execution

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass: `go test ./...`
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.