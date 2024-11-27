package cmd

import (
	"github.com/spf13/cobra"
)

var (
	accountName       string
	showBlockingUsers bool
	showBlockList     bool
	showListCount     bool
	showListNames     bool
	rootCmd           = &cobra.Command{
		Use:   "clearsky",
		Short: "clearsky is a command line interface for the clearsky.app API",
		Long:  "clearsky is a command line interface for the clearsky.app API",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&accountName, "account", "a", "", "The BlueSky account name.")
}
