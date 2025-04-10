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

// ListInstances creates a tool to list instances within an environment
func ListInstances(getClient GetClientFn) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_instances",
			mcp.WithDescription("List all instances within a KB Cloud environment"),
			mcp.WithString("org_name",
				mcp.Required(),
				mcp.Description("Organization name"),
			),
			mcp.WithString("env_name",
				mcp.Required(),
				mcp.Description("Environment name"),
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

			// Call KB Cloud API - using ClusterApi's ListInstance method
			instances, resp, err := client.Cluster.ListInstance(client.Context, orgName, envName)
			if err != nil {
				return nil, fmt.Errorf("failed to list instances: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			// Check response status
			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to list instances: %s", string(body))), nil
			}

			// Return result
			result, err := json.Marshal(instances)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(result)), nil
		}
}

// GetInstance creates a tool to get details of a specific instance
func GetInstance(getClient GetClientFn) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_instance",
			mcp.WithDescription("Get details of a specific instance in KB Cloud"),
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

			// Call KB Cloud API - using ClusterApi's GetCluster method
			// Note: In KB Cloud API, instances are referred to as clusters
			instance, resp, err := client.Cluster.GetCluster(client.Context, orgName, envName, instanceName)
			if err != nil {
				return nil, fmt.Errorf("failed to get instance: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			// Check response status
			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get instance: %s", string(body))), nil
			}

			// Return result
			result, err := json.Marshal(instance)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(result)), nil
		}
}
