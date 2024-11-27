package cmd

import (
	"github.com/spf13/cobra"
)

var (
	accountName       string
	err               error
	showBlockingUsers bool
	// showListNames     bool
	userId  string
	rootCmd = &cobra.Command{
		Use:   "clearsky",
		Short: "clearsky is a command line interface for the clearsky.app API",
		Long:  "clearsky is a command line interface for the clearsky.app API",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

// func init() {

// 	return
// }
