package cmd

import (
	"fmt"
	"os"

	"github.com/gdanko/clearsky/globals"
	"github.com/gdanko/clearsky/pkg/api"
	"github.com/gdanko/clearsky/util"
	"github.com/markkurossi/tabulate"
	"github.com/spf13/cobra"
)

var (
	blocksCmd = &cobra.Command{
		Use:          "blocks",
		Short:        "Display the number of users blocking --account.",
		Long:         "Display the number of users blocking --account.",
		PreRunE:      blocksPreRunCmd,
		RunE:         blocksRunCmd,
		SilenceUsage: true,
	}
)

func init() {
	GetBlocksFlags(blocksCmd)
	rootCmd.AddCommand(blocksCmd)
}

func blocksPreRunCmd(cmd *cobra.Command, args []string) error {
	globals.SetDebugFlag(debugFlag)
	if accountName != "" {
		displayName, userId, err = api.GetUserID(accountName)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("The required --account flag is missing")
		cmd.Help()
		os.Exit(1)
	}
	return nil
}

func blocksRunCmd(cmd *cobra.Command, args []string) error {
	var (
		blockListOutput    globals.BlockListOutput
		chunk              []globals.BlockingUser
		divided            [][]globals.BlockingUser
		err                error
		totalRecords       int
		newBlockListOutput globals.BlockListOutput
	)
	blockListOutput, err = api.GetBlockingUsersList(userId)
	if err != nil {
		return err
	}

	// https://medium.com/insiderengineering/concurrent-http-requests-in-golang-best-practices-and-techniques-f667e5a19dea
	totalRecords = len(blockListOutput.Items)
	if listMaxResults < totalRecords {
		blockListOutput.Items = blockListOutput.Items[0:listMaxResults]
	}
	if showBlockingUsers {
		alignment := tabulate.ML
		tab := tabulate.New(
			tabulate.Unicode,
		)
		tab.Header("id").SetAlign(alignment)
		tab.Header("handle").SetAlign(alignment)
		tab.Header("name").SetAlign(alignment)
		divided = util.SliceChunker(blockListOutput.Items, batchChunkSize)
		for _, chunk = range divided {
			api.ExpandBlockListUsers(&chunk, batchOperationTimeout)
			newBlockListOutput.Items = append(newBlockListOutput.Items, chunk...)
		}
		for _, item := range newBlockListOutput.Items {
			row := tab.Row()
			row.Column(item.DID)
			row.Column(item.Username)
			row.Column(item.DisplayName)
		}
		tab.Print(os.Stdout)
	}
	if listMaxResults < totalRecords {
		fmt.Printf("Results limited to %d entries by use of --limit.\n", listMaxResults)
	}
	fmt.Printf("%s (%s) is currently being blocked by %d users\n", accountName, displayName, len(blockListOutput.Items))

	return nil
}
