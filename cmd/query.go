package cmd

import (
	"context"
	"strings"

	"github.com/mikesmitty/file-search/internal/constants"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:     "query [text...]",
	Aliases: []string{"q"},
	Short:   "Query Gemini File Search",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := getClient(ctx)
		if err != nil {
			return err
		}
		defer client.Close()

		// Resolve store name to ID if --store was used
		storeID := queryStoreID
		if queryStoreName != "" {
			storeID, err = client.ResolveStoreName(ctx, queryStoreName)
			if err != nil {
				return err
			}
		}

		if queryModel == "" {
			queryModel = constants.DefaultModel
		}

		// Join all arguments to form the query string
		queryString := strings.Join(args, " ")

		resp, err := client.Query(ctx, queryString, storeID, queryModel, queryMetadataFilter)
		if err != nil {
			return err
		}
		return printOutput(resp, outputFormat)
	},
}

var (
	queryStoreName      string
	queryStoreID        string
	queryModel          string
	queryMetadataFilter string
)

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.Flags().StringVar(&queryStoreName, "store", "", "Store display name (optional)")
	queryCmd.Flags().StringVar(&queryStoreID, "store-id", "", "Store resource ID (optional, "+constants.StoreResourcePrefix+"xxx)")
	queryCmd.Flags().StringVar(&queryModel, "model", constants.DefaultModel, "Model name")
	queryCmd.Flags().StringVar(&queryMetadataFilter, "metadata-filter", "", "Metadata filter expression (optional)")
	queryCmd.RegisterFlagCompletionFunc("store", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getCompleter().GetStoreNames(), cobra.ShellCompDirectiveNoFileComp
	})
	queryCmd.RegisterFlagCompletionFunc("store-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getCompleter().GetStoreNames(), cobra.ShellCompDirectiveNoFileComp
	})
	queryCmd.RegisterFlagCompletionFunc("model", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getCompleter().GetModelNames(), cobra.ShellCompDirectiveNoFileComp
	})
}
