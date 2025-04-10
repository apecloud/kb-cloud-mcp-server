package kbcloud

import (
	"github.com/apecloud/kb-cloud-mcp-server/pkg/translations"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTools registers all KB Cloud MCP tools with the MCP server
func RegisterTools(s *server.MCPServer, getClientFn GetClientFn, t translations.TranslationHelperFunc) {
	// Organization tools
	organizationTool, organizationHandler := ListOrganizations(getClientFn)
	s.AddTool(organizationTool, organizationHandler)

	orgDetailTool, orgDetailHandler := GetOrganization(getClientFn)
	s.AddTool(orgDetailTool, orgDetailHandler)

	// Environment tools
	environmentsTool, environmentsHandler := ListEnvironments(getClientFn)
	s.AddTool(environmentsTool, environmentsHandler)

	envDetailTool, envDetailHandler := GetEnvironment(getClientFn)
	s.AddTool(envDetailTool, envDetailHandler)

	// Instance tools
	instancesTool, instancesHandler := ListInstances(getClientFn)
	s.AddTool(instancesTool, instancesHandler)

	instanceDetailTool, instanceDetailHandler := GetInstance(getClientFn)
	s.AddTool(instanceDetailTool, instanceDetailHandler)

	// Backup tools
	backupsTool, backupsHandler := ListBackups(getClientFn)
	s.AddTool(backupsTool, backupsHandler)

	backupDetailTool, backupDetailHandler := GetBackup(getClientFn)
	s.AddTool(backupDetailTool, backupDetailHandler)
}
