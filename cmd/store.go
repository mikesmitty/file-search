package cmd

import (
	"context"
	"fmt"

	"github.com/mikesmitty/file-search/internal/constants"
	"github.com/mikesmitty/file-search/internal/gemini"
	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Manage File Search Stores",
}

func init() {
	rootCmd.AddCommand(storeCmd)

	// Store list
	storeCmd.AddCommand(&cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all File Search Stores",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := getClient(ctx)
			if err != nil {
				return err
			}
			defer client.Close()
			stores, err := client.ListStores(ctx)
			if err != nil {
				return err
			}
			return printOutput(stores, outputFormat)
		},
	})

	// Store get
	storeCmd.AddCommand(&cobra.Command{
		Use:   "get [name]",
		Short: "Get details of a File Search Store",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return getCompleter().GetStoreNames(), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := getClient(ctx)
			if err != nil {
				return err
			}
			defer client.Close()

			// Resolve store name to ID
			storeID, err := client.ResolveStoreName(ctx, args[0])
			if err != nil {
				return err
			}

			store, err := client.GetStore(ctx, storeID)
			if err != nil {
				return err
			}
			return printOutput(store, outputFormat)
		},
	})

	// Store delete
	var deleteStoreForce bool
	deleteStoreCmd := &cobra.Command{
		Use:     "delete [name]",
		Aliases: []string{"rm", "del"},
		Short:   "Delete a File Search Store",
		Args:    cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return getCompleter().GetStoreNames(), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := getClient(ctx)
			if err != nil {
				return err
			}
			defer client.Close()

			// Resolve store name to ID
			storeID, err := client.ResolveStoreName(ctx, args[0])
			if err != nil {
				return err
			}

			err = client.DeleteStore(ctx, storeID, deleteStoreForce)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printOutput(map[string]string{"status": "deleted", "name": args[0]}, "json")
			}
			fmt.Printf("Deleted store: %s\n", args[0])
			return nil
		},
	}
	deleteStoreCmd.Flags().BoolVar(&deleteStoreForce, "force", false, "Force delete even if store contains documents")
	storeCmd.AddCommand(deleteStoreCmd)

	// Store create
	storeCmd.AddCommand(&cobra.Command{
		Use:     "create [display_name]",
		Aliases: []string{"new", "add"},
		Short:   "Create a new File Search Store",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := getClient(ctx)
			if err != nil {
				return err
			}
			defer client.Close()
			store, err := client.CreateStore(ctx, args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printOutput(store, "json")
			}
			fmt.Printf("Created store: %s (%s)\n", store.DisplayName, store.Name)
			return nil
		},
	})

	// Store import-file
	var importFileStore string
	var importFileStoreID string
	importFileCmd := &cobra.Command{
		Use:   "import-file [file-name-or-id]",
		Short: "Import a file from Files API into a Store",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if importFileStore == "" && importFileStoreID == "" {
				return fmt.Errorf("either --store or --store-id is required")
			}
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

			// Resolve store name to ID if --store was used
			storeID := importFileStoreID
			if importFileStore != "" {
				storeID, err = client.ResolveStoreName(ctx, importFileStore)
				if err != nil {
					return err
				}
			}

			err = client.ImportFile(ctx, fileID, storeID, &gemini.ImportFileOptions{
				Quiet: quiet,
			})
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printOutput(map[string]string{"status": "imported", "file": fileID, "store": storeID}, "json")
			}
			// ImportFile already prints progress if not quiet, but we can add a final success message if needed
			// The client.ImportFile method prints "Import complete" so we are good.
			return nil
		},
	}
	importFileCmd.Flags().StringVar(&importFileStore, "store", "", "Store display name")
	importFileCmd.Flags().StringVar(&importFileStoreID, "store-id", "", "Store resource ID ("+constants.StoreResourcePrefix+"xxx)")
	importFileCmd.RegisterFlagCompletionFunc("store", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getCompleter().GetStoreNames(), cobra.ShellCompDirectiveNoFileComp
	})
	importFileCmd.RegisterFlagCompletionFunc("store-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getCompleter().GetStoreNames(), cobra.ShellCompDirectiveNoFileComp
	})
	storeCmd.AddCommand(importFileCmd)
}
