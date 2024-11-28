package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/gdanko/clearsky/globals"
	"github.com/gdanko/clearsky/pkg/api"
	"github.com/spf13/cobra"
)

var (
	countsCmd = &cobra.Command{
		Use:          "counts",
		Short:        "Show active users, deleted users, and total users.",
		Long:         "Show active users, deleted users, and total users.",
		PreRun:       countsPreRunCmd,
		Run:          countsRunCmd,
		SilenceUsage: true,
	}
)

func init() {
	rootCmd.AddCommand(countsCmd)
}

func countsPreRunCmd(cmd *cobra.Command, args []string) {
	// Empty
}

func countsRunCmd(cmd *cobra.Command, args []string) {
	var (
		countData globals.CountData
		body      []byte
		url       string
	)

	url = "https://api.clearsky.services/api/v1/anon/total-users"
	body, err = api.FetchUrl(url, logger)
	if err != nil {
		panic(err)
	}
	countData = globals.CountData{}
	err = json.Unmarshal(body, &countData)
	if err != nil {
		panic(err)
	}
	fmt.Printf("As of: %s\n", countData.Data.AsOf)
	fmt.Printf("  Total users:   %s\n", countData.Data.TotalCount.Value)
	fmt.Printf("  Deleted users: %s\n", countData.Data.DeletedCount.Value)
	fmt.Printf("  Active users:  %s\n\n", countData.Data.ActiveCount.Value)
}
