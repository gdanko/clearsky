package cmd

import (
	"fmt"
	"os"

	"github.com/gdanko/clearsky/globals"
	"github.com/gdanko/clearsky/pkg/api"
	"github.com/spf13/cobra"
)

var (
	listsCmd = &cobra.Command{
		Use:          "lists",
		Short:        "Display the number of moderated lists --account is a member of.",
		Long:         "Display the number of moderated lists --account is a member of.",
		PreRunE:      listsPreRunCmd,
		RunE:         listsRunCmd,
		SilenceUsage: true,
	}
)

func init() {
	GetListsFlags(listsCmd)
	rootCmd.AddCommand(listsCmd)
}

func listsPreRunCmd(cmd *cobra.Command, args []string) error {
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

func listsRunCmd(cmd *cobra.Command, args []string) error {
	return nil
}
