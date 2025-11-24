package cmd

import (
	"context"

	"github.com/mikesmitty/file-search/internal/gemini"
	"github.com/mikesmitty/file-search/internal/mcp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP Server",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		// For MCP, we start the server even without API key configured.
		// Tools will fail gracefully when invoked if auth is missing.
		key, _ := getAPIKey()

		var client mcp.GeminiClient
		if key != "" {
			c, err := gemini.NewClient(ctx, key, nil)
			if err != nil {
				return err
			}
			defer c.Close()
			client = c
		}

		tools := getMCPTools()
		return mcp.RunServer(ctx, client, tools)
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)

	mcpCmd.Flags().StringVar(&mcpTools, "mcp-tools", "", "Comma-separated list of MCP tools to enable (default: query)")
	viper.BindPFlag("mcp_tools", mcpCmd.Flags().Lookup("mcp-tools"))
	viper.BindEnv("mcp_tools", "MCP_TOOLS")
}
