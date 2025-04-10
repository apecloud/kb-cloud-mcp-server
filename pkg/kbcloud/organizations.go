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

// ListOrganizations creates a tool to list organizations in KB Cloud
func ListOrganizations(getClient GetClientFn) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_organizations",
			mcp.WithDescription("List all organizations you have access to in KB Cloud"),
			WithPagination(),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Get pagination parameters
			pagination, err := OptionalPaginationParams(request)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			// Get KB Cloud client
			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get KB Cloud client: %w", err)
			}

			// Call KB Cloud API
			opts := make([]interface{}, 0)
			orgs, resp, err := client.Organization.ListOrg(client.Context, opts...)
			if err != nil {
				return nil, fmt.Errorf("failed to list organizations: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			// Check response status
			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to list organizations: %s", string(body))), nil
			}

			// Return result
			result, err := json.Marshal(orgs)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(result)), nil
		}
}

// GetOrganization creates a tool to get details of a specific organization
func GetOrganization(getClient GetClientFn) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_organization",
			mcp.WithDescription("Get details of a specific organization in KB Cloud"),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Organization name"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Get required parameters
			orgName, err := RequiredParam[string](request, "name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			// Get KB Cloud client
			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get KB Cloud client: %w", err)
			}

			// Call KB Cloud API
			org, resp, err := client.Organization.ReadOrg(client.Context, orgName)
			if err != nil {
				return nil, fmt.Errorf("failed to get organization: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			// Check response status
			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get organization: %s", string(body))), nil
			}

			// Return result
			result, err := json.Marshal(org)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(result)), nil
		}
}
