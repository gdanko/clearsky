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
	blockingCmd = &cobra.Command{
		Use:          "blocking",
		Short:        "Display the number of users --account is blocking.",
		Long:         "Display the number of users --account is blocking.",
		PreRunE:      blockingPreRunCmd,
		RunE:         blockingRunCmd,
		SilenceUsage: true,
	}
)

func init() {
	GetBlockingFlags(blockingCmd)
	rootCmd.AddCommand(blockingCmd)
}

func blockingPreRunCmd(cmd *cobra.Command, args []string) error {
	logLevel = logLevelMap[logLevelStr]
	logger = util.ConfigureLogger(logLevel, nocolorFlag)

	if accountName != "" {
		// Get the target's DID
		targetDid, err = api.GetUserDid(accountName, logger)
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

func blockingRunCmd(cmd *cobra.Command, args []string) error {
	var (
		alignment   tabulate.Align
		blockedList map[string]globals.BlockingUser
		err         error
		item        globals.BlockingUser
		row         *tabulate.Row
		tab         *tabulate.Tabulate
	)
	blockedList, err = api.GetBlocking(userId, showBlockedUsers, batchOperationTimeout, listMaxResults, logger)
	if err != nil {
		return err
	}
	if showBlockedUsers {
		alignment = tabulate.ML
		tab = tabulate.New(
			tabulate.Unicode,
		)
		tab.Header("DID").SetAlign(alignment)
		tab.Header("Handle").SetAlign(alignment)
		tab.Header("Display Name").SetAlign(alignment)

		for _, item = range blockedList {
			row = tab.Row()
			row.Column(item.DID)
			row.Column(item.Handle)
			row.Column(util.StripNonPrintable(item.DisplayName))
		}
		tab.Print(os.Stdout)
	}
	if displayName == "" {
		fmt.Printf("%s is currently blocking %s users\n", accountName, util.AddCommas(len(blockedList)))
	} else {
		fmt.Printf("%s (%s) is currently blocking %s users\n", accountName, displayName, util.AddCommas(len(blockedList)))
	}

	return nil
}
