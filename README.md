# KubeBlocks Cloud MCP Server

The KubeBlocks Cloud MCP Server is a [Model Context Protocol (MCP)](https://modelcontextprotocol.io/introduction)
server that provides seamless integration with KubeBlocks Cloud APIs, enabling AI assistants to interact with
KubeBlocks Cloud resources through a standardized tool-calling interface.

[![Go Version](https://img.shields.io/github/go-mod/go-version/apecloud/kb-cloud-mcp-server?style=flat-square)](https://go.mod)
[![License](https://img.shields.io/github/license/apecloud/kb-cloud-mcp-server?style=flat-square)](LICENSE)

## Use Cases

- Automating KubeBlocks Cloud resource management
- Retrieving and analyzing data from KubeBlocks Cloud environments
- Building AI-powered tools that interact with the KubeBlocks Cloud ecosystem
- Enabling AI assistants to provision and manage database instances

## Features

- MCP-based API for accessing KubeBlocks Cloud resources
- Secure authentication via KubeBlocks Cloud API key and secret
- Support for common KubeBlocks Cloud resources:
  - Organizations
  - Environments
  - Instances
  - Backups
- Internationalization support via translation helpers
- Secure communication through StdioServer

## Prerequisites

1. Go 1.20+ (compatible with the MCP-Go package)
2. KubeBlocks Cloud API credentials - you'll need an API key name and secret

## Installation

### Build from Source

```bash
git clone https://github.com/apecloud/kb-cloud-mcp-server.git
cd kb-cloud-mcp-server
go mod tidy
go build -o kb-cloud-mcp-server ./cmd/server
```

### Usage with VS Code

Add the following JSON block to your User Settings (JSON) file in VS Code. You can do this by pressing `Ctrl + Shift + P` and typing `Preferences: Open User Settings (JSON)`.

Optionally, you can add it to a file called `.vscode/mcp.json` in your workspace to share the configuration with others.

> Note that the `mcp` key is not needed in the `.vscode/mcp.json` file.

```json
{
  "mcp": {
    "inputs": [
      {
        "type": "promptString",
        "id": "kb_cloud_api_key",
        "description": "KubeBlocks Cloud API Key Name",
        "password": false
      },
      {
        "type": "promptString",
        "id": "kb_cloud_api_secret",
        "description": "KubeBlocks Cloud API Secret",
        "password": true
      }
    ],
    "servers": {
      "kbcloud": {
        "command": "docker",
        "args": [
          "run",
          "-i",
          "--rm",
          "-e",
          "KB_CLOUD_API_KEY_NAME",
          "-e",
          "KB_CLOUD_API_KEY_SECRET",
          "apecloud/kb-cloud-mcp-server:latest"
        ],
        "env": {
          "KB_CLOUD_API_KEY_NAME": "${input:kb_cloud_api_key}",
          "KB_CLOUD_API_KEY_SECRET": "${input:kb_cloud_api_secret}"
        }
      }
    }
  }
}
```

More about using MCP server tools in VS Code's [agent mode documentation](https://code.visualstudio.com/docs/copilot/chat/mcp-servers).

## Usage

### Environment Variables

The server supports the following environment variables:

```bash
# KubeBlocks Cloud API credentials (required)
export KB_CLOUD_API_KEY_NAME=your-api-key-name
export KB_CLOUD_API_KEY_SECRET=your-api-key-secret
export KB_CLOUD_SITE=https://api.apecloud.com  # Optional: KubeBlocks Cloud API endpoint

# Server configuration
export KB_CLOUD_MCP_LOG_LEVEL=info  # debug, info, warn, error
export KB_CLOUD_MCP_EXPORT_TRANSLATIONS=true  # Optional: Export translations to JSON file
export KB_CLOUD_DEBUG=true  # Optional: Enable debug mode for KubeBlocks Cloud API client
```

### Starting the Server

```bash
./kb-cloud-mcp-server
```

Or with command-line flags:

```bash
./kb-cloud-mcp-server stdio --api-key=your-api-key-name --api-secret=your-api-key-secret
```

### Configuration File

You can also use a configuration file:

```yaml
# .kb-cloud-mcp-server.yaml
log_level: info
api_key: your-api-key-name
api_secret: your-api-key-secret
site_url: https://api.apecloud.com
```

Then start the server with:

```bash
./kb-cloud-mcp-server stdio --config=.kb-cloud-mcp-server.yaml
```

## Available MCP Tools

The server provides the following MCP tools for interacting with KubeBlocks Cloud resources:

### Organizations

- **list_organizations** - List all organizations you have access to
  - No parameters required

- **get_organization** - Get details of a specific organization
  - `organizationId`: Organization unique identifier (string, required)

### Environments

- **list_environments** - List all environments within an organization
  - `organizationId`: Organization unique identifier (string, required)

- **get_environment** - Get details of a specific environment
  - `organizationId`: Organization unique identifier (string, required)
  - `environmentId`: Environment unique identifier (string, required)

### Instances

- **list_instances** - List all instances within an environment
  - `organizationId`: Organization unique identifier (string, required)
  - `environmentId`: Environment unique identifier (string, required)

- **get_instance** - Get details of a specific instance
  - `organizationId`: Organization unique identifier (string, required)
  - `environmentId`: Environment unique identifier (string, required)
  - `instanceId`: Instance unique identifier (string, required)

### Backups

- **list_backups** - List all backups for an instance
  - `organizationId`: Organization unique identifier (string, required)
  - `environmentId`: Environment unique identifier (string, required)
  - `instanceId`: Instance unique identifier (string, required)

- **get_backup** - Get details of a specific backup
  - `organizationId`: Organization unique identifier (string, required)
  - `environmentId`: Environment unique identifier (string, required)
  - `instanceId`: Instance unique identifier (string, required)
  - `backupId`: Backup unique identifier (string, required)

## Library Usage

The exported Go API of this module should currently be considered unstable and subject to breaking changes. In the future, we may offer stability; please file an issue if there is a use case where this would be valuable.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](./LICENSE) file for details.
