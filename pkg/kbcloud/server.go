package kbcloud

import (
	"os"

	"github.com/apecloud/kb-cloud-mcp-server/pkg/translations"
	"github.com/mark3labs/mcp-go/server"
)

// NewServer creates a new KB Cloud MCP server
func NewServer(version string) *server.MCPServer {
	// Initialize translation helper
	t, dumpTranslations := translations.TranslationHelper()

	// Create a new MCP server
	s := server.NewMCPServer(
		"kb-cloud-mcp-server",
		version,
		server.WithLogging(),
	)

	// Register KB Cloud tools
	getClientFn := GetDefaultClientFn()
	RegisterTools(s, getClientFn, t)

	// Export translations if requested
	// Get environment variable using os.LookupEnv directly to avoid conflict
	if val, exists := os.LookupEnv("KB_CLOUD_MCP_EXPORT_TRANSLATIONS"); exists && (val == "true" || val == "1") {
		dumpTranslations()
	}

	return s
}
