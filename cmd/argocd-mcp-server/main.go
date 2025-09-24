package main

import (
	"context"
	"log"
	"os"

	"argo_mcp/internal/server"
)

func main() {
	// Set up logging
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Create context
	ctx := context.Background()

	// Create and start the MCP server
	mcpServer := server.NewMCPServer()

	log.Println("Starting MCP server...")
	if err := mcpServer.Run(ctx); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}