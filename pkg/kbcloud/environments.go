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

// ListEnvironments creates a tool to list environments within an organization
func ListEnvironments(getClient GetClientFn) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_environments",
			mcp.WithDescription("List all environments within a KB Cloud organization"),
			mcp.WithString("org_name",
				mcp.Required(),
				mcp.Description("Organization name"),
			),
			WithPagination(),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Get required parameters
			orgName, err := RequiredParam[string](request, "org_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

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
			envs, resp, err := client.Environment.ListEnvironment(client.Context, orgName)
			if err != nil {
				return nil, fmt.Errorf("failed to list environments: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			// Check response status
			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to list environments: %s", string(body))), nil
			}

			// Return result
			result, err := json.Marshal(envs)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(result)), nil
		}
}

// GetEnvironment creates a tool to get details of a specific environment
func GetEnvironment(getClient GetClientFn) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_environment",
			mcp.WithDescription("Get details of a specific environment in KB Cloud"),
			mcp.WithString("org_name",
				mcp.Required(),
				mcp.Description("Organization name"),
			),
			mcp.WithString("env_name",
				mcp.Required(),
				mcp.Description("Environment name"),
			),
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

			// Get KB Cloud client
			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get KB Cloud client: %w", err)
			}

			// Call KB Cloud API
			env, resp, err := client.Environment.GetEnvironment(client.Context, orgName, envName)
			if err != nil {
				return nil, fmt.Errorf("failed to get environment: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			// Check response status
			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get environment: %s", string(body))), nil
			}

			// Return result
			result, err := json.Marshal(env)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(result)), nil
		}
}
