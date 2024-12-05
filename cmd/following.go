package cmd

import (
	"fmt"
	"os"

	"github.com/gdanko/clearsky/pkg/api"
	"github.com/gdanko/clearsky/util"
	"github.com/spf13/cobra"
)

var (
	followingCmd = &cobra.Command{
		Use:          "following",
		Short:        "Display the number users --account is following.",
		Long:         "Display the number users --account is following.",
		PreRunE:      followingPreRunCmd,
		RunE:         followingRunCmd,
		SilenceUsage: true,
	}
)

func init() {
	GetFollowingFlags(followingCmd)
	rootCmd.AddCommand(followingCmd)
}

func followingPreRunCmd(cmd *cobra.Command, args []string) error {
	logLevel = logLevelMap[logLevelStr]
	logger = util.ConfigureLogger(logLevel, nocolorFlag)

	if accountName != "" {
		targetDid, err = api.GetUserDid(accountName, logger)
		if err != nil {
			return err
		}
		targetCredentials, err = api.GetTargetInfo(targetDid, logger)
		// Validate the target info
	} else {
		fmt.Println("The required --account flag is missing")
		cmd.Help()
		os.Exit(1)
	}
	return nil
}

func followingRunCmd(cmd *cobra.Command, args []string) error {
	var (
		// alignment     tabulate.Align
		followingList map[string]api.Following2
		err           error
		// item          api.Following
		// row           *tabulate.Row
		// tab           *tabulate.Tabulate
	)
	followingList, err = api.GetFollowing2(targetCredentials, showFollowingUsers, batchOperationTimeout, listMaxResults, logger)
	if err != nil {
		return err
	}
	// if showFollowingUsers {
	// 	alignment = tabulate.ML
	// 	tab = tabulate.New(
	// 		tabulate.Unicode,
	// 	)
	// 	tab.Header("DID").SetAlign(alignment)
	// 	tab.Header("Handle").SetAlign(alignment)
	// 	tab.Header("Display Name").SetAlign(alignment)

	// 	for _, item = range followingList {
	// 		row = tab.Row()
	// 		row.Column(item.DID)
	// 		row.Column(item.Handle)
	// 		row.Column(util.StripNonPrintable(item.DisplayName))
	// 	}
	// 	tab.Print(os.Stdout)
	// }
	if displayName == "" {
		fmt.Printf("%s is currently following %s users\n", accountName, util.AddCommas(len(followingList)))
	} else {
		fmt.Printf("%s (%s) is currently following %s users\n", accountName, displayName, util.AddCommas(len(followingList)))
	}

	return nil
}
