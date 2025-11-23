package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mikesmitty/file-search/internal/constants"
	"github.com/mikesmitty/file-search/internal/gemini"
	"github.com/spf13/cobra"
)

var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "Manage Files",
}

func init() {
	rootCmd.AddCommand(fileCmd)

	// File list
	fileCmd.AddCommand(&cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List uploaded files",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := getClient(ctx)
			if err != nil {
				return err
			}
			defer client.Close()
			files, err := client.ListFiles(ctx)
			if err != nil {
				return err
			}
			return printOutput(files, outputFormat)
		},
	})

	// File get
	fileCmd.AddCommand(&cobra.Command{
		Use:   "get [name]",
		Short: "Get details of a file",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return getCompleter().GetFileNames(), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := getClient(ctx)
			if err != nil {
				return err
			}
			defer client.Close()

			// Resolve file name to ID
			fileID, err := client.ResolveFileName(ctx, args[0])
			if err != nil {
				return err
			}

			file, err := client.GetFile(ctx, fileID)
			if err != nil {
				return err
			}
			return printOutput(file, outputFormat)
		},
	})

	// File delete
	fileCmd.AddCommand(&cobra.Command{
		Use:     "delete [name]",
		Aliases: []string{"rm", "del"},
		Short:   "Delete a file",
		Args:    cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return getCompleter().GetFileNames(), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := getClient(ctx)
			if err != nil {
				return err
			}
			defer client.Close()

			// Resolve file name to ID
			fileID, err := client.ResolveFileName(ctx, args[0])
			if err != nil {
				return err
			}

			err = client.DeleteFile(ctx, fileID)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printOutput(map[string]string{"status": "deleted", "file": fileID}, "json")
			}
			fmt.Printf("Deleted file: %s\n", args[0])
			return nil
		},
	})

	// File upload
	var uploadStoreName string
	var uploadStoreID string
	var uploadDisplayName string
	var uploadMimeType string
	var uploadChunkSize int
	var uploadChunkOverlap int
	var uploadMetadata []string
	uploadCmd := &cobra.Command{
		Use:   "upload [path]",
		Short: "Upload and import a file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := getClient(ctx)
			if err != nil {
				return err
			}
			defer client.Close()

			// Parse metadata from key=value strings
			metadataMap := make(map[string]string)
			for _, meta := range uploadMetadata {
				parts := strings.SplitN(meta, "=", 2)
				if len(parts) == 2 {
					metadataMap[parts[0]] = parts[1]
				}
			}

			// Auto-set display name from filename if not provided
			displayName := uploadDisplayName
			if displayName == "" {
				displayName = filepath.Base(args[0])
			}

			// Resolve store name to ID if --store was used
			storeID := uploadStoreID
			if uploadStoreName != "" {
				storeID, err = client.ResolveStoreName(ctx, uploadStoreName)
				if err != nil {
					return err
				}
			}

			opts := &gemini.UploadFileOptions{
				StoreName:      storeID,
				DisplayName:    displayName,
				MIMEType:       uploadMimeType,
				MaxChunkTokens: uploadChunkSize,
				ChunkOverlap:   uploadChunkOverlap,
				Metadata:       metadataMap,
				Quiet:          quiet,
			}

			file, err := client.UploadFile(ctx, args[0], opts)
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				if file != nil {
					return printOutput(file, "json")
				}
				return printOutput(map[string]string{"status": "uploaded_and_indexed"}, "json")
			}

			if file != nil {
				fmt.Printf("Uploaded file: %s (URI: %s)\n", file.DisplayName, file.URI)
			}
			return nil
		},
	}
	uploadCmd.Flags().StringVar(&uploadStoreName, "store", "", "Store display name (optional)")
	uploadCmd.Flags().StringVar(&uploadStoreID, "store-id", "", "Store resource ID (optional, "+constants.StoreResourcePrefix+"xxx)")
	uploadCmd.Flags().StringVar(&uploadDisplayName, "name", "", "Display name (optional)")
	uploadCmd.Flags().StringVar(&uploadMimeType, "mime-type", "", "MIME type (optional, e.g. text/plain, application/pdf)")
	uploadCmd.Flags().IntVar(&uploadChunkSize, "chunk-size", 0, "Max tokens per chunk (for store uploads)")
	uploadCmd.Flags().IntVar(&uploadChunkOverlap, "chunk-overlap", 0, "Overlap tokens between chunks (for store uploads)")
	uploadCmd.Flags().StringArrayVar(&uploadMetadata, "metadata", []string{}, "Custom metadata as key=value (repeatable, for store uploads)")
	uploadCmd.RegisterFlagCompletionFunc("store", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getCompleter().GetStoreNames(), cobra.ShellCompDirectiveNoFileComp
	})
	uploadCmd.RegisterFlagCompletionFunc("store-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getCompleter().GetStoreNames(), cobra.ShellCompDirectiveNoFileComp
	})
	fileCmd.AddCommand(uploadCmd)
}
