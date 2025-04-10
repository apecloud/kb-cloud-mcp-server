package kbcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ListBackups creates a tool to list backups for an instance
func ListBackups(getClient GetClientFn) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_backups",
			mcp.WithDescription("List all backups for a KB Cloud instance"),
			mcp.WithString("org_name",
				mcp.Required(),
				mcp.Description("Organization name"),
			),
			mcp.WithString("env_name",
				mcp.Required(),
				mcp.Description("Environment name"),
			),
			mcp.WithString("instance_name",
				mcp.Required(),
				mcp.Description("Instance name"),
			),
			WithPagination(),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Get required parameters
			orgName, err := RequiredParam[string](request, "org_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			envName, err := RequiredParam[string](request, "env_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			instanceName, err := RequiredParam[string](request, "instance_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			// Get KB Cloud client
			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get KB Cloud client: %w", err)
			}

			// Call KB Cloud API
			backups, resp, err := client.Backup.ListBackups(client.Context, orgName, &[]interface{}{})
			if err != nil {
				return nil, fmt.Errorf("failed to list backups: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			// Check response status
			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to list backups: %s", string(body))), nil
			}

			// Return result
			result, err := json.Marshal(backups)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(result)), nil
		}
}

// GetBackup creates a tool to get details of a specific backup
func GetBackup(getClient GetClientFn) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_backup",
			mcp.WithDescription("Get details of a specific backup in KB Cloud"),
			mcp.WithString("org_name",
				mcp.Required(),
				mcp.Description("Organization name"),
			),
			mcp.WithString("backup_id",
				mcp.Required(),
				mcp.Description("Backup ID"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Get required parameters
			orgName, err := RequiredParam[string](request, "org_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			backupID, err := RequiredParam[string](request, "backup_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			// Get KB Cloud client
			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get KB Cloud client: %w", err)
			}

			// Call KB Cloud API
			backup, resp, err := client.Backup.GetBackup(client.Context, orgName, backupID)
			if err != nil {
				return nil, fmt.Errorf("failed to get backup: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			// Check response status
			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get backup: %s", string(body))), nil
			}

			// Return result
			result, err := json.Marshal(backup)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(result)), nil
		}
}
