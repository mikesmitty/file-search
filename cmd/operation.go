package cmd

import (
	"context"
	"fmt"

	"github.com/mikesmitty/file-search/internal/gemini"
	"github.com/spf13/cobra"
)

var operationCmd = &cobra.Command{
	Use:     "operation",
	Aliases: []string{"op", "ops", "operations"},
	Short:   "Manage long-running operations",
}

func init() {
	rootCmd.AddCommand(operationCmd)

	var operationType string
	operationGetCmd := &cobra.Command{
		Use:   "get [operation-name]",
		Short: "Get the status of a long-running operation",
		Long: `Get the status of a long-running file upload or import operation.

Operation names follow the format: fileSearchStores/{store-id}/operations/{operation-id}

Examples:
  # Get operation status (auto-detect type)
  file-search operation get "fileSearchStores/abc123/operations/op456"

  # Get operation status with specific type
  file-search operation get "fileSearchStores/abc123/operations/op456" --type import

  # Get operation status in JSON format
  file-search operation get "fileSearchStores/abc123/operations/op456" --format json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := getClient(ctx)
			if err != nil {
				return err
			}
			defer client.Close()

			var opType gemini.OperationType
			switch operationType {
			case "import":
				opType = gemini.OperationTypeImport
			case "upload":
				opType = gemini.OperationTypeUpload
			case "":
				// Auto-detect (empty string is valid)
				opType = ""
			default:
				return fmt.Errorf("invalid operation type: %s (must be 'import' or 'upload')", operationType)
			}

			status, err := client.GetOperation(ctx, args[0], opType)
			if err != nil {
				return err
			}

			return printOutput(status, outputFormat)
		},
	}
	operationGetCmd.Flags().StringVar(&operationType, "type", "", "Operation type: import or upload (auto-detect if not specified)")
	operationCmd.AddCommand(operationGetCmd)
}
