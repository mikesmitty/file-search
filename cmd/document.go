package cmd

import (
	"context"
	"fmt"

	"github.com/mikesmitty/file-search/internal/constants"
	"github.com/spf13/cobra"
)

var documentCmd = &cobra.Command{
	Use:     "document",
	Aliases: []string{"doc", "docs"},
	Short:   "Manage Documents in Stores",
}

func init() {
	rootCmd.AddCommand(documentCmd)

	// Document list
	var docListStore string
	var docListStoreID string
	docListCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List documents in a store",
		RunE: func(cmd *cobra.Command, args []string) error {
			if docListStore == "" && docListStoreID == "" {
				return fmt.Errorf("either --store or --store-id is required")
			}
			ctx := context.Background()
			client, err := getClient(ctx)
			if err != nil {
				return err
			}
			defer client.Close()

			// Resolve store name to ID if --store was used
			storeID := docListStoreID
			if docListStore != "" {
				storeID, err = client.ResolveStoreName(ctx, docListStore)
				if err != nil {
					return err
				}
			}

			docs, err := client.ListDocuments(ctx, storeID)
			if err != nil {
				return err
			}
			return printOutput(docs, outputFormat)
		},
	}
	docListCmd.Flags().StringVar(&docListStore, "store", "", "Store display name")
	docListCmd.Flags().StringVar(&docListStoreID, "store-id", "", "Store resource ID ("+constants.StoreResourcePrefix+"xxx)")
	docListCmd.RegisterFlagCompletionFunc("store", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getCompleter().GetStoreNames(), cobra.ShellCompDirectiveNoFileComp
	})
	docListCmd.RegisterFlagCompletionFunc("store-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getCompleter().GetStoreNames(), cobra.ShellCompDirectiveNoFileComp
	})
	documentCmd.AddCommand(docListCmd)

	// Document get
	var docGetStore string
	var docGetStoreID string
	docGetCmd := &cobra.Command{
		Use:   "get [name]",
		Short: "Get document details",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			storeFlag, _ := cmd.Flags().GetString("store")
			if storeFlag != "" {
				return getCompleter().GetDocumentNames(storeFlag), cobra.ShellCompDirectiveNoFileComp
			}
			storeIDFlag, _ := cmd.Flags().GetString("store-id")
			if storeIDFlag != "" {
				return getCompleter().GetDocumentNames(storeIDFlag), cobra.ShellCompDirectiveNoFileComp
			}
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := getClient(ctx)
			if err != nil {
				return err
			}
			defer client.Close()

			// If store is provided, resolve document name within that store
			docID := args[0]
			if docGetStore != "" || docGetStoreID != "" {
				storeRef := docGetStoreID
				if docGetStore != "" {
					storeRef = docGetStore
				}
				docID, err = client.ResolveDocumentName(ctx, storeRef, args[0])
				if err != nil {
					return err
				}
			}

			doc, err := client.GetDocument(ctx, docID)
			if err != nil {
				return err
			}
			return printOutput(doc, outputFormat)
		},
	}
	docGetCmd.Flags().StringVar(&docGetStore, "store", "", "Store display name (optional, for name resolution)")
	docGetCmd.Flags().StringVar(&docGetStoreID, "store-id", "", "Store resource ID (optional, for name resolution)")
	docGetCmd.RegisterFlagCompletionFunc("store", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getCompleter().GetStoreNames(), cobra.ShellCompDirectiveNoFileComp
	})
	docGetCmd.RegisterFlagCompletionFunc("store-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getCompleter().GetStoreNames(), cobra.ShellCompDirectiveNoFileComp
	})
	documentCmd.AddCommand(docGetCmd)

	// Document delete
	var docDelStore string
	var docDelStoreID string
	var docDelForce bool
	docDelCmd := &cobra.Command{
		Use:     "delete [name]",
		Aliases: []string{"rm", "del"},
		Short:   "Delete a document",
		Args:    cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			storeFlag, _ := cmd.Flags().GetString("store")
			if storeFlag != "" {
				return getCompleter().GetDocumentNames(storeFlag), cobra.ShellCompDirectiveNoFileComp
			}
			storeIDFlag, _ := cmd.Flags().GetString("store-id")
			if storeIDFlag != "" {
				return getCompleter().GetDocumentNames(storeIDFlag), cobra.ShellCompDirectiveNoFileComp
			}
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := getClient(ctx)
			if err != nil {
				return err
			}
			defer client.Close()

			// If store is provided, resolve document name within that store
			docID := args[0]
			if docDelStore != "" || docDelStoreID != "" {
				storeRef := docDelStoreID
				if docDelStore != "" {
					storeRef = docDelStore
				}
				docID, err = client.ResolveDocumentName(ctx, storeRef, args[0])
				if err != nil {
					return err
				}
			}

			err = client.DeleteDocument(ctx, docID, docDelForce)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printOutput(map[string]string{"status": "deleted", "document": docID}, "json")
			}
			fmt.Printf("Deleted document: %s\n", args[0])
			return nil
		},
	}
	docDelCmd.Flags().StringVar(&docDelStore, "store", "", "Store display name (optional, for name resolution)")
	docDelCmd.Flags().StringVar(&docDelStoreID, "store-id", "", "Store resource ID (optional, for name resolution)")
	docDelCmd.Flags().BoolVar(&docDelForce, "force", false, "Force delete even if document contains chunks")
	docDelCmd.RegisterFlagCompletionFunc("store", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getCompleter().GetStoreNames(), cobra.ShellCompDirectiveNoFileComp
	})
	docDelCmd.RegisterFlagCompletionFunc("store-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getCompleter().GetStoreNames(), cobra.ShellCompDirectiveNoFileComp
	})
	documentCmd.AddCommand(docDelCmd)
}
