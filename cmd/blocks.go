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
	logLevel = logLevelMap[logLevelStr]
	logger = util.ConfigureLogger(logLevel, nocolorFlag)

	if accountName != "" {
		displayName, userId, err = api.GetUserID(accountName, logger)
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
		blockingList map[string]globals.BlockingUser
		err          error
	)

	blockingList, err = api.GetBlockingUsersList(userId, batchChunkSize, batchOperationTimeout, logger)
	if err != nil {
		return err
	}

	if showBlockingUsers {
		alignment := tabulate.ML
		tab := tabulate.New(
			tabulate.Unicode,
		)
		tab.Header("DID").SetAlign(alignment)
		tab.Header("Handle").SetAlign(alignment)
		tab.Header("Display Name").SetAlign(alignment)

		for _, item := range blockingList {
			row := tab.Row()
			row.Column(item.DID)
			row.Column(item.Username)
			row.Column(util.StripNonPrintable(item.DisplayName))
		}
		tab.Print(os.Stdout)
	}
	fmt.Printf("%s (%s) is currently being blocked by %d users\n", accountName, displayName, len(blockingList))

	return nil
}
