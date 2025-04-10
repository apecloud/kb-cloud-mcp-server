package kbcloud

import (
	"context"
	"fmt"
	"os"

	"github.com/apecloud/kb-cloud-client-go/api/common"
	"github.com/apecloud/kb-cloud-client-go/api/kbcloud"
)

// GetClientFn is a function type that returns a KB Cloud API client
type GetClientFn func(ctx context.Context) (*Client, error)

// Client is a wrapper for the KB Cloud API client
type Client struct {
	APIClient *common.APIClient
	Context   context.Context

	// Organization API
	Organization *kbcloud.OrganizationApi

	// Environment API
	Environment *kbcloud.EnvironmentApi

	// Cluster API (for instances)
	Cluster *kbcloud.ClusterApi

	// Backup API
	Backup *kbcloud.BackupApi
}

// NewClient creates a new KB Cloud client
func NewClient(apiClient *common.APIClient, ctx context.Context) *Client {
	return &Client{
		APIClient:    apiClient,
		Context:      ctx,
		Organization: kbcloud.NewOrganizationApi(apiClient),
		Environment:  kbcloud.NewEnvironmentApi(apiClient),
		Cluster:      kbcloud.NewClusterApi(apiClient),
		Backup:       kbcloud.NewBackupApi(apiClient),
	}
}

// GetDefaultClientFn returns a function that creates a KB Cloud client from request context
func GetDefaultClientFn() GetClientFn {
	return func(ctx context.Context) (*Client, error) {
		// Extract API key and secret from the context
		apiKey, apiSecret, ok := GetAPICredentials(ctx)
		if !ok {
			return nil, fmt.Errorf("KB Cloud API credentials not found in context")
		}

		// Set authentication context for KB Cloud API
		authCtx := context.WithValue(
			ctx,
			common.ContextDigestAuth,
			common.DigestAuth{
				UserName: apiKey,
				Password: apiSecret,
			},
		)

		// Get site configuration if provided
		site, hasSite := GetSiteConfiguration(ctx)
		if hasSite {
			authCtx = context.WithValue(
				authCtx,
				common.ContextServerVariables,
				map[string]string{"site": site},
			)
		}

		// Create configuration
		config := common.NewConfiguration()

		// Set debug mode based on context
		if isDebug(ctx) {
			config.Debug = true
		}

		// Create API client
		apiClient := common.NewAPIClient(config)

		// Create and return the KB Cloud client
		return NewClient(apiClient, authCtx), nil
	}
}

// GetAPICredentials extracts the KB Cloud API credentials from the context
func GetAPICredentials(ctx context.Context) (apiKey, apiSecret string, ok bool) {
	// Try to get from context values
	toolContext, ok := ctx.Value("toolContext").(map[string]string)
	if ok {
		if apiKey, ok := toolContext["KB_CLOUD_API_KEY_NAME"]; ok {
			if apiSecret, ok := toolContext["KB_CLOUD_API_KEY_SECRET"]; ok {
				return apiKey, apiSecret, true
			}
		}
	}

	// Fallback to environment variables
	apiKey, apiSecret = getEnvAPICredentials()
	return apiKey, apiSecret, apiKey != "" && apiSecret != ""
}

// GetSiteConfiguration extracts the KB Cloud site configuration from the context
func GetSiteConfiguration(ctx context.Context) (site string, ok bool) {
	// Try to get from context values
	toolContext, ok := ctx.Value("toolContext").(map[string]string)
	if ok {
		if site, ok := toolContext["KB_CLOUD_SITE"]; ok && site != "" {
			return site, true
		}
	}

	// Fallback to environment variables
	site = getEnvSiteConfiguration()
	return site, site != ""
}

// getEnvAPICredentials gets API credentials from environment variables
func getEnvAPICredentials() (apiKey, apiSecret string) {
	apiKey = getEnvOrDefault("KB_CLOUD_API_KEY_NAME", "")
	apiSecret = getEnvOrDefault("KB_CLOUD_API_KEY_SECRET", "")
	return apiKey, apiSecret
}

// getEnvSiteConfiguration gets site configuration from environment variables
func getEnvSiteConfiguration() string {
	return getEnvOrDefault("KB_CLOUD_SITE", "")
}

// getEnvOrDefault gets an environment variable or returns the default value
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := getEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnv gets an environment variable
func getEnv(key string) (string, bool) {
	val, ok := os.LookupEnv(key)
	return val, ok
}

// isDebug checks if debug mode is enabled in the context
func isDebug(ctx context.Context) bool {
	toolContext, ok := ctx.Value("toolContext").(map[string]string)
	if ok {
		if debug, ok := toolContext["KB_CLOUD_DEBUG"]; ok {
			return debug == "true" || debug == "1"
		}
	}

	// Fallback to environment variable
	debug, _ := getEnv("KB_CLOUD_DEBUG")
	return debug == "true" || debug == "1"
}
