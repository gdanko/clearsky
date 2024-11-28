package cmd

import "github.com/spf13/cobra"

func GetBlocksFlags(cmd *cobra.Command) {
	getBlocksAndListsFlags(cmd)
	getBlocksFlags(cmd)
}

func GetListsFlags(cmd *cobra.Command) {
	getBlocksAndListsFlags(cmd)
	getListsFlags(cmd)
}

func getBlocksAndListsFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&accountName, "account", "a", "", "The BlueSky account name.")
	cmd.Flags().IntVarP(&listMaxResults, "limit", "l", 0, "Limit the results to --limit - for testing.")
	cmd.Flags().IntVarP(&batchOperationTimeout, "timeout", "t", 60, "When making batched http requests, specify the timeout in seconds.")
	cmd.Flags().IntVarP(&batchChunkSize, "chunk-size", "c", 20, "When making batched http requests, specify the chunk size.")
}

func getBlocksFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&showBlockingUsers, "blocking-users", "u", false, "Gather the list of blocking users' names (expensive).")
}

func getListsFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&showBlockingUsers, "list-names", "n", false, "Gather the list of moderated lists' names (expensive).")
}
