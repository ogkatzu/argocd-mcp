package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

// TestArgocdConnection tests the connection from MCP server to ArgoCD
func TestArgocdConnection(t *testing.T) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		t.Logf("Warning: Could not load .env file: %v", err)
	}

	// Test configuration loading
	t.Run("Configuration", func(t *testing.T) {
		testConfiguration(t)
	})

	// Test network connectivity
	t.Run("NetworkConnectivity", func(t *testing.T) {
		testNetworkConnectivity(t)
	})

	// Test authentication
	t.Run("Authentication", func(t *testing.T) {
		testAuthentication(t)
	})

	// Test API endpoints
	t.Run("APIEndpoints", func(t *testing.T) {
		testAPIEndpoints(t)
	})
}

func testConfiguration(t *testing.T) {
	serverURL := getEnvWithDefault("ARGOCD_SERVER", "https://localhost:8080")
	authToken := os.Getenv("ARGOCD_AUTH_TOKEN")
	insecure := getEnvWithDefault("ARGOCD_INSECURE", "true") == "true"

	t.Logf("ArgoCD Server URL: %s", serverURL)
	t.Logf("Auth Token Present: %t", authToken != "")
	t.Logf("TLS Skip Verify: %t", insecure)

	if serverURL == "" {
		t.Error("ARGOCD_SERVER is not set")
	}

	if authToken == "" || authToken == "your-token-here" {
		t.Error("ARGOCD_AUTH_TOKEN is not set or using default placeholder value")
	}
}

func testNetworkConnectivity(t *testing.T) {
	serverURL := getEnvWithDefault("ARGOCD_SERVER", "https://localhost:8080")
	insecure := getEnvWithDefault("ARGOCD_INSECURE", "true") == "true"

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test basic connectivity to the server
	req, err := http.NewRequestWithContext(ctx, "GET", serverURL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Network connectivity test failed: %v", err)
	}
	defer resp.Body.Close()

	t.Logf("Server responded with status: %d", resp.StatusCode)

	// Any response (even 4xx/5xx) indicates network connectivity is working
	if resp.StatusCode == 0 {
		t.Error("No response received from server")
	}
}

func testAuthentication(t *testing.T) {
	serverURL := getEnvWithDefault("ARGOCD_SERVER", "https://localhost:8080")
	authToken := os.Getenv("ARGOCD_AUTH_TOKEN")
	insecure := getEnvWithDefault("ARGOCD_INSECURE", "true") == "true"

	if authToken == "" || authToken == "your-token-here" {
		t.Skip("Skipping authentication test: no valid auth token provided")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test authentication with a simple endpoint
	url := fmt.Sprintf("%s/api/v1/session", serverURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Authentication test request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Auth test response status: %d", resp.StatusCode)
	t.Logf("Auth test response body: %s", string(body))

	if resp.StatusCode == 401 {
		t.Error("Authentication failed: Invalid or expired token")
	} else if resp.StatusCode == 403 {
		t.Error("Authentication failed: Insufficient permissions")
	} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		t.Log("Authentication successful")
	}
}

func testAPIEndpoints(t *testing.T) {
	serverURL := getEnvWithDefault("ARGOCD_SERVER", "https://localhost:8080")
	authToken := os.Getenv("ARGOCD_AUTH_TOKEN")
	insecure := getEnvWithDefault("ARGOCD_INSECURE", "true") == "true"

	if authToken == "" || authToken == "your-token-here" {
		t.Skip("Skipping API endpoint test: no valid auth token provided")
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test the applications endpoint (what the MCP server uses)
	url := fmt.Sprintf("%s/api/v1/applications", serverURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Applications API test request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	t.Logf("Applications API response status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Applications API failed with status %d: %s", resp.StatusCode, string(body))
		return
	}

	// Try to parse the response
	var appList struct {
		Items []map[string]interface{} `json:"items"`
	}

	if err := json.Unmarshal(body, &appList); err != nil {
		t.Errorf("Failed to parse applications response: %v", err)
		t.Logf("Raw response: %s", string(body))
		return
	}

	t.Logf("Successfully retrieved %d applications", len(appList.Items))

	// Log first few applications for debugging
	for i, app := range appList.Items {
		if i >= 3 { // Only show first 3
			break
		}
		if metadata, ok := app["metadata"].(map[string]interface{}); ok {
			if name, ok := metadata["name"].(string); ok {
				t.Logf("Application %d: %s", i+1, name)
			}
		}
	}
}

// TestArgocdConnectionVerbose runs a verbose version that prints all details
func TestArgocdConnectionVerbose(t *testing.T) {
	fmt.Println("=== ArgoCD MCP Server Connection Test ===")
	fmt.Println()

	// Load environment
	fmt.Println("Loading environment configuration...")
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Could not load .env file: %v\n", err)
	}

	serverURL := getEnvWithDefault("ARGOCD_SERVER", "https://localhost:8080")
	authToken := os.Getenv("ARGOCD_AUTH_TOKEN")
	insecure := getEnvWithDefault("ARGOCD_INSECURE", "true") == "true"

	fmt.Printf("Configuration:\n")
	fmt.Printf("  Server URL: %s\n", serverURL)
	fmt.Printf("  Auth Token: %s\n", maskToken(authToken))
	fmt.Printf("  Skip TLS Verify: %t\n", insecure)
	fmt.Println()

	// Create HTTP client
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		},
	}

	// Test network connectivity
	fmt.Println("Testing network connectivity...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", serverURL, nil)
	if err != nil {
		fmt.Printf("❌ Failed to create request: %v\n", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Network connectivity failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("✅ Server reachable (status: %d)\n", resp.StatusCode)
	fmt.Println()

	// Test authentication
	if authToken == "" || authToken == "your-token-here" {
		fmt.Println("⚠️  No valid auth token configured - skipping auth test")
		return
	}

	fmt.Println("Testing authentication...")
	url := fmt.Sprintf("%s/api/v1/applications", serverURL)
	req, err = http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Printf("❌ Failed to create auth request: %v\n", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("❌ Auth test failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 401 {
		fmt.Printf("❌ Authentication failed: Invalid or expired token\n")
		fmt.Printf("Response: %s\n", string(body))
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ API call failed with status %d: %s\n", resp.StatusCode, string(body))
		return
	}

	fmt.Printf("✅ Authentication successful\n")

	// Parse applications
	var appList struct {
		Items []map[string]interface{} `json:"items"`
	}

	if err := json.Unmarshal(body, &appList); err != nil {
		fmt.Printf("❌ Failed to parse applications: %v\n", err)
		return
	}

	fmt.Printf("✅ Successfully retrieved %d applications\n", len(appList.Items))
	fmt.Println()

	fmt.Println("=== Connection Test Complete ===")
}

// Helper functions
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func maskToken(token string) string {
	if token == "" {
		return "(not set)"
	}
	if token == "your-token-here" {
		return "(placeholder - needs actual token)"
	}
	if len(token) < 8 {
		return "***"
	}
	return token[:4] + "..." + token[len(token)-4:]
}