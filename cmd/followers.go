package cmd

import (
	"fmt"
	"os"

	"github.com/gdanko/clearsky/pkg/api"
	"github.com/gdanko/clearsky/util"
	"github.com/markkurossi/tabulate"
	"github.com/spf13/cobra"
)

var (
	followersCmd = &cobra.Command{
		Use:          "followers",
		Short:        "Display the number of followers --account has.",
		Long:         "Display the number of followers --account has.",
		PreRunE:      followersPreRunCmd,
		RunE:         followersRunCmd,
		SilenceUsage: true,
	}
)

func init() {
	GetFollowersFlags(followersCmd)
	rootCmd.AddCommand(followersCmd)
}

func followersPreRunCmd(cmd *cobra.Command, args []string) error {
	logLevel = logLevelMap[logLevelStr]
	logger = util.ConfigureLogger(logLevel, nocolorFlag)

	if accountName != "" {
		targetDid, err = api.GetUserDid(accountName, logger)
		if err != nil {
			return err
		}
		targetCredentials, err = api.GetTargetInfo(targetDid, logger)
	} else {
		fmt.Println("The required --account flag is missing")
		cmd.Help()
		os.Exit(1)
	}
	return nil
}

func followersRunCmd(cmd *cobra.Command, args []string) error {
	var (
		alignment     tabulate.Align
		followersList map[string]api.Follower
		err           error
		item          api.Follower
		row           *tabulate.Row
		tab           *tabulate.Tabulate
	)
	followersList, err = api.GetFollowers(targetCredentials, showFollowingUsers, batchOperationTimeout, listMaxResults, logger)
	if err != nil {
		return err
	}
	if showFollowingUsers {
		alignment = tabulate.ML
		tab = tabulate.New(
			tabulate.Unicode,
		)
		tab.Header("DID").SetAlign(alignment)
		tab.Header("Handle").SetAlign(alignment)
		tab.Header("Display Name").SetAlign(alignment)

		for _, item = range followersList {
			row = tab.Row()
			row.Column(item.DID)
			row.Column(item.Handle)
			row.Column(util.StripNonPrintable(item.DisplayName))
		}
		tab.Print(os.Stdout)
	}
	if displayName == "" {
		fmt.Printf("%s is currently being followed by %s users\n", accountName, util.AddCommas(len(followersList)))
	} else {
		fmt.Printf("%s (%s) is currently being followed by %s users\n", accountName, displayName, util.AddCommas(len(followersList)))
	}

	return nil
}
