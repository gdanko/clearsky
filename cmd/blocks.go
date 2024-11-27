package cmd

import "github.com/spf13/cobra"

var (
	blocksCmd = &cobra.Command{
		Use:          "blocks",
		Short:        "Display the number of users blocking --account",
		Long:         "Display the number of users blocking --account",
		PreRun:       blocksPreRunCmd,
		Run:          blocksRunCmd,
		SilenceUsage: true,
	}
)

func init() {
	blocksCmd.PersistentFlags().BoolVarP(&showBlockingUsers, "blocking-users", "u", false, "Gather the list of blocking users' names (expensive).")
	rootCmd.AddCommand(blocksCmd)
}

func blocksPreRunCmd(cmd *cobra.Command, args []string) {
	// Empty
}

func blocksRunCmd(cmd *cobra.Command, args []string) {
	// Empty
}
