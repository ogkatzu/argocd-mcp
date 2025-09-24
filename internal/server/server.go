package server

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPServer represents our ArgoCD MCP server instance
type MCPServer struct {
	server     *mcp.Server
	config     *ServerConfig
	status     *ServerStatus
	argocdCfg  *ArgocdConfig
	httpClient *http.Client
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

// ArgocdConfig holds ArgoCD connection configuration
type ArgocdConfig struct {
	ServerURL   string `json:"server_url"`
	AuthToken   string `json:"auth_token,omitempty"`
	Insecure    bool   `json:"insecure"`
}

// ArgocdApplication represents an ArgoCD application
type ArgocdApplication struct {
	Metadata struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
	Spec struct {
		Project     string `json:"project"`
		Source      struct {
			RepoURL        string `json:"repoURL"`
			Path           string `json:"path"`
			TargetRevision string `json:"targetRevision"`
		} `json:"source"`
		Destination struct {
			Server    string `json:"server"`
			Namespace string `json:"namespace"`
		} `json:"destination"`
	} `json:"spec"`
	Status struct {
		Sync struct {
			Status string `json:"status"`
		} `json:"sync"`
		Health struct {
			Status string `json:"status"`
		} `json:"health"`
	} `json:"status"`
}

// ArgocdApplicationList represents a list of ArgoCD applications
type ArgocdApplicationList struct {
	Items []ArgocdApplication `json:"items"`
}

// ServerStatus holds server runtime status
type ServerStatus struct {
	StartTime    time.Time `json:"start_time"`
	RequestCount int64     `json:"request_count"`
	LastRequest  time.Time `json:"last_request"`
}

// NewMCPServer creates a new ArgoCD MCP server instance
func NewMCPServer() *MCPServer {
	// Load .env file if it exists (non-fatal if it doesn't)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading .env: %v", err)
	}

	config := &ServerConfig{
		Name:        "argocd-mcp-server",
		Version:     "1.0.0",
		Description: "ArgoCD MCP server for managing GitOps deployments",
	}

	status := &ServerStatus{
		StartTime: time.Now(),
	}

	// Initialize ArgoCD configuration from environment variables
	argocdCfg := &ArgocdConfig{
		ServerURL: getEnvWithDefault("ARGOCD_SERVER", "https://localhost:8080"),
		AuthToken: os.Getenv("ARGOCD_AUTH_TOKEN"),
		Insecure:  getEnvWithDefault("ARGOCD_INSECURE", "true") == "true",
	}


	// Create HTTP client with optional TLS skip
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: argocdCfg.Insecure,
			},
		},
	}

	mcpServer := &MCPServer{
		config:     config,
		status:     status,
		argocdCfg:  argocdCfg,
		httpClient: httpClient,
	}

	// Create the MCP server with implementation info
	impl := &mcp.Implementation{
		Name:    config.Name,
		Version: config.Version,
	}

	server := mcp.NewServer(impl, nil)

	mcpServer.server = server
	mcpServer.setupHandlers()

	return mcpServer
}

// setupHandlers configures all the MCP handlers
func (s *MCPServer) setupHandlers() {
	// TODO: Add ArgoCD-specific tools here
	// Examples:
	// - list_applications
	// - get_application_status
	// - sync_application
	// - create_application
	// - delete_application
	// - get_cluster_info
	// - etc.

	// Add ArgoCD applications resource
	s.server.AddResource(&mcp.Resource{
		URI:         "argocd://applications",
		Name:        "ArgoCD Applications",
		Description: "List of all ArgoCD applications",
		MIMEType:    "application/json",
	}, s.handleApplicationsResource)
}

// Run starts the ArgoCD MCP server
func (s *MCPServer) Run(ctx context.Context) error {
	log.Printf("Starting %s v%s", s.config.Name, s.config.Version)
	log.Printf("Server description: %s", s.config.Description)

	// Run the server using stdio transport
	return s.server.Run(ctx, &mcp.StdioTransport{})
}

// Resource handlers

func (s *MCPServer) handleApplicationsResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	s.updateRequestStats()

	// Make API call to ArgoCD
	apps, err := s.getArgocdApplications(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get ArgoCD applications: %w", err)
	}

	// Convert to JSON
	appsJSON, err := json.MarshalIndent(apps, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal applications: %w", err)
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "argocd://applications",
				MIMEType: "application/json",
				Text:     string(appsJSON),
			},
		},
	}, nil
}

// ArgoCD API functions

func (s *MCPServer) getArgocdApplications(ctx context.Context) (*ArgocdApplicationList, error) {
	url := fmt.Sprintf("%s/api/v1/applications", s.argocdCfg.ServerURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header if token is available
	if s.argocdCfg.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+s.argocdCfg.AuthToken)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ArgoCD API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var appList ArgocdApplicationList
	if err := json.Unmarshal(body, &appList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &appList, nil
}

// Helper functions

func (s *MCPServer) updateRequestStats() {
	s.status.RequestCount++
	s.status.LastRequest = time.Now()
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}