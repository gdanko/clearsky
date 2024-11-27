package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/gdanko/clearsky/globals"
	"github.com/gdanko/clearsky/pkg/api"
	"github.com/gdanko/clearsky/util"
	"github.com/kr/pretty"
	"github.com/spf13/cobra"
)

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
	blocksCmd.PersistentFlags().StringVarP(&accountName, "account", "a", "", "The BlueSky account name.")
	blocksCmd.PersistentFlags().BoolVarP(&showBlockingUsers, "blocking-users", "u", false, "Gather the list of blocking users' names (expensive).")
	rootCmd.AddCommand(blocksCmd)
}

func blocksPreRunCmd(cmd *cobra.Command, args []string) {
	userId, err = api.GetUserID(accountName)
	if err != nil {
		panic(err)
	}
}

func blocksRunCmd(cmd *cobra.Command, args []string) {
	var (
		batchOperationTimeout = 60
		blockListOutput       globals.BlockListOutput
		blockListPage         globals.BlockListPage
		body                  []byte
		chunk                 []globals.BlockingUser
		chunkSize             = 35
		divided               [][]globals.BlockingUser
		i                     int
		maxPages              = 1000
		newBlockListOutput    globals.BlockListOutput
		url                   string
	)
	url = fmt.Sprintf("https://api.clearsky.services/api/v1/anon/single-blocklist/%s", userId)
	body, err = api.FetchUrl(url)
	if err != nil {
		panic(err)
	}

	blockListPage = globals.BlockListPage{}
	err = json.Unmarshal(body, &blockListPage)
	if err != nil {
		panic(err)
	}
	if len(blockListPage.Data.Blocklist) > 0 {
		blockListOutput.Items = append(blockListOutput.Items, blockListPage.Data.Blocklist...)
	} else {
		panic(err)
	}

	// Now we cycle through /2, /3, etc until there are no more
	if len(blockListOutput.Items) >= 100 {
		for i := 2; i <= maxPages; i++ {
			url = fmt.Sprintf("https://api.clearsky.services/api/v1/anon/single-blocklist/%s/%d", userId, i)
			body, err = api.FetchUrl(url)
			if err != nil {
				panic(err)
			}
			blockListPage = globals.BlockListPage{}
			err = json.Unmarshal(body, &blockListPage)
			if err != nil {
				panic(err)
			}
			if len(blockListPage.Data.Blocklist) > 0 {
				blockListOutput.Items = append(blockListOutput.Items, blockListPage.Data.Blocklist...)
			} else {
				break
			}
		}
	}

	// https://medium.com/insiderengineering/concurrent-http-requests-in-golang-best-practices-and-techniques-f667e5a19dea
	// blockListOutput.Items = blockListOutput.Items[0:100]
	fmt.Println(len(blockListOutput.Items))
	if showBlockingUsers {
		divided = util.SliceChunker(blockListOutput.Items, chunkSize)
		for i, chunk = range divided {
			fmt.Printf("Chunk %d\n", i)
			api.ExpandBlockListUsers(&chunk, batchOperationTimeout)
			newBlockListOutput.Items = append(newBlockListOutput.Items, chunk...)
			// fmt.Printf("Sleeping for %d seconds\n", sleepSeconds)
			// time.Sleep(time.Duration(sleepSeconds) * time.Second)
		}
	}
	newBlockListOutput.Count = len(newBlockListOutput.Items)
	pretty.Println(newBlockListOutput.Items)
	fmt.Printf("%s is currently being blocked by %d users\n", accountName, newBlockListOutput.Count)

}
