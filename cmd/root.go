package cmd

import (
	"github.com/spf13/cobra"
)

var (
	accountName           string
	batchOperationTimeout int
	batchChunkSize        int
	debugFlag             bool
	err                   error
	listMaxResults        int
	showBlockingUsers     bool
	showListNames         bool
	displayName           string
	userId                string
	rootCmd               = &cobra.Command{
		Use:   "clearsky",
		Short: "clearsky is a command line interface for the clearsky.services API.",
		Long:  "clearsky is a command line interface for the clearsky.services API.",
	}
)

func Execute() error {
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "d", false, "Enable debugging output.")

	return rootCmd.Execute()
}
