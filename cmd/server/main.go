package main

import (
	"context"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/apecloud/kb-cloud-mcp-server/pkg/kbcloud"
	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "version"
var commit = "commit"
var date = "date"

var (
	rootCmd = &cobra.Command{
		Use:     "server",
		Short:   "KB Cloud MCP Server",
		Long:    `KB Cloud MCP Server provides a Model Context Protocol (MCP) server for KubeBlocks Cloud.`,
		Version: fmt.Sprintf("%s (%s) %s", version, commit, date),
	}

	stdioCmd = &cobra.Command{
		Use:   "stdio",
		Short: "Start stdio server",
		Long:  `Start a server that communicates via standard input/output streams using JSON-RPC messages.`,
		Run: func(_ *cobra.Command, _ []string) {
			logFile := viper.GetString("log-file")
			logger, err := initLogger(logFile)
			if err != nil {
				stdlog.Fatal("Failed to initialize logger:", err)
			}

			apiKey := viper.GetString("api-key")
			apiSecret := viper.GetString("api-secret")
			siteURL := viper.GetString("site-url")

			cfg := runConfig{
				logger:    logger,
				apiKey:    apiKey,
				apiSecret: apiSecret,
				siteURL:   siteURL,
			}

			if err := runStdioServer(cfg); err != nil {
				stdlog.Fatal("failed to run stdio server:", err)
			}
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	// Add global flags
	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.kb-cloud-mcp-server.yaml)")
	rootCmd.PersistentFlags().String("log-file", "", "Path to log file")
	rootCmd.PersistentFlags().String("log-level", "info", "Log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().String("api-key", "", "KB Cloud API key name")
	rootCmd.PersistentFlags().String("api-secret", "", "KB Cloud API key secret")
	rootCmd.PersistentFlags().String("site-url", "", "KB Cloud site URL")

	// Bind to viper
	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("log-file", rootCmd.PersistentFlags().Lookup("log-file"))
	_ = viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	_ = viper.BindPFlag("api-key", rootCmd.PersistentFlags().Lookup("api-key"))
	_ = viper.BindPFlag("api-secret", rootCmd.PersistentFlags().Lookup("api-secret"))
	_ = viper.BindPFlag("site-url", rootCmd.PersistentFlags().Lookup("site-url"))

	// Add subcommands
	rootCmd.AddCommand(stdioCmd)
}

func initConfig() {
	// Handle config file if provided
	configFile := viper.GetString("config")
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		// Look for config in current directory with name .kb-cloud-mcp-server
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME")
		viper.SetConfigName(".kb-cloud-mcp-server")
	}

	// Set environment variable prefix
	viper.SetEnvPrefix("KB_CLOUD_MCP")

	// Replace '-' with '_' in env vars
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)

	// Read from environment variables
	viper.AutomaticEnv()

	// Read config if present
	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	}

	// Override with environment variables
	if v := os.Getenv("KB_CLOUD_API_KEY_NAME"); v != "" {
		viper.Set("api-key", v)
	}
	if v := os.Getenv("KB_CLOUD_API_KEY_SECRET"); v != "" {
		viper.Set("api-secret", v)
	}
	if v := os.Getenv("KB_CLOUD_SITE"); v != "" {
		viper.Set("site-url", v)
	}
}

func initLogger(outPath string) (*log.Logger, error) {
	logger := log.New()

	// Set log level
	level, err := log.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		level = log.InfoLevel
	}
	logger.SetLevel(level)

	if outPath == "" {
		logger.SetFormatter(&log.JSONFormatter{})
		return logger, nil
	}

	file, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	logger.SetOutput(file)
	logger.SetFormatter(&log.JSONFormatter{})

	return logger, nil
}

type runConfig struct {
	logger    *log.Logger
	apiKey    string
	apiSecret string
	siteURL   string
}

func runStdioServer(cfg runConfig) error {
	// Create app context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Set environment variables for KB Cloud client
	if cfg.apiKey != "" {
		os.Setenv("KB_CLOUD_API_KEY_NAME", cfg.apiKey)
	}
	if cfg.apiSecret != "" {
		os.Setenv("KB_CLOUD_API_KEY_SECRET", cfg.apiSecret)
	}
	if cfg.siteURL != "" {
		os.Setenv("KB_CLOUD_SITE", cfg.siteURL)
	}

	// Create MCP server
	s := server.NewMCPServer(
		"kb-cloud-mcp-server",
		version,
		server.WithLogging(),
	)

	// Register KB Cloud tools
	kbcloud.RegisterTools(s)

	// Create stdio server
	stdioServer := server.NewStdioServer(s)
	stdLogger := stdlog.New(cfg.logger.Writer(), "stdioserver", 0)
	stdioServer.SetErrorLogger(stdLogger)

	// Start listening for messages
	errC := make(chan error, 1)
	go func() {
		in, out := io.Reader(os.Stdin), io.Writer(os.Stdout)
		errC <- stdioServer.Listen(ctx, in, out)
	}()

	// Output server running message
	_, _ = fmt.Fprintf(os.Stderr, "KB Cloud MCP Server running on stdio\n")

	// Wait for shutdown signal
	select {
	case <-ctx.Done():
		cfg.logger.Info("Shutting down server...")
	case err := <-errC:
		if err != nil {
			return fmt.Errorf("error running server: %w", err)
		}
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
