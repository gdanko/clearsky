package cmd

import "github.com/spf13/cobra"

var (
	listsCmd = &cobra.Command{
		Use:          "lists",
		Short:        "Display the number of moderated lists --account is a member of",
		Long:         "Display the number of moderated lists --account is a member of",
		PreRun:       listsPreRunCmd,
		Run:          listsRunCmd,
		SilenceUsage: true,
	}
)

func init() {
	listsCmd.PersistentFlags().BoolVarP(&showBlockingUsers, "list-names", "n", false, "Gather the list of moderated lists' names (expensive).")
	rootCmd.AddCommand(listsCmd)
}

func listsPreRunCmd(cmd *cobra.Command, args []string) {
	// Empty
}

func listsRunCmd(cmd *cobra.Command, args []string) {
	// Empty
}
