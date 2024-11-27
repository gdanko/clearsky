package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/gdanko/clearsky/globals"
	"github.com/gdanko/clearsky/pkg/api"
	"github.com/spf13/cobra"
)

var (
	accountName       string
	showBlockingUsers bool
	showBlockList     bool
	showListCount     bool
	showListNames     bool
	userId            string
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

	url := fmt.Sprintf("https://api.clearsky.services/api/v1/anon/get-did/%s", accountName)
	body, err := api.FetchUrl(url)
	if err != nil {
		panic(err)
	}

	getDid := globals.GetUserData()
	err = json.Unmarshal(body, &getDid)
	if err != nil {
		panic(err)
	}
	userId = getDid.Data.DidIdentifier
}
