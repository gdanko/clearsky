package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/forPelevin/gomoji"
	"github.com/gdanko/clearsky/globals"
	"github.com/sirupsen/logrus"
	"github.com/useinsider/go-pkg/insrequester"
	// "golang.org/x/sync/errgroup"
)

func GetUserID(accountName string, logger *logrus.Logger) (displayName string, userId string, err error) {
	var (
		body []byte
		url  string
	)
	// Get userId from handle
	url = fmt.Sprintf("https://api.clearsky.services/api/v1/anon/get-did/%s", accountName)
	body, err = FetchUrl(url, logger)
	if err != nil {
		return displayName, userId, err
	}

	getDid := globals.UserDid{}
	err = json.Unmarshal(body, &getDid)
	if err != nil {
		return displayName, userId, err
	}
	userId = getDid.Data.DidIdentifier

	// Get displayName from userId
	url = fmt.Sprintf("https://public.api.bsky.app/xrpc/app.bsky.actor.getProfile?actor=%s", (userId))
	body, err = FetchUrl(url, logger)
	if err != nil {
		return displayName, userId, err
	}

	userInfo := globals.BlueSkyUser{}
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return displayName, userId, err
	}
	displayName = userInfo.DisplayName

	return displayName, userId, nil
}

func worker(requester *insrequester.Request, jobs <-chan globals.Job, results chan<- *http.Response, wg *sync.WaitGroup, logger *logrus.Logger) {
	for job := range jobs {
		logger.Debug(job.URL)
		res, err := requester.Get(insrequester.RequestEntity{Endpoint: job.URL})
		if err != nil {
			fmt.Println(err)
		}
		results <- res
		wg.Done()
	}
}

func GetBlockingUsersList(userId string, logger *logrus.Logger) (blockListOutput globals.BlockListOutput, err error) {
	var (
		blockListPage globals.BlockListPage
		body          []byte
		maxPages      = 1000
		url           string
	)

	// https://api.clearsky.services/api/v1/anon/single-blocklist/did:plc:ccskhvd467uwdrxpwaudnbni
	url = fmt.Sprintf("https://api.clearsky.services/api/v1/anon/single-blocklist/%s", userId)
	body, err = FetchUrl(url, logger)
	if err != nil {
		return globals.BlockListOutput{}, err
	}

	blockListPage = globals.BlockListPage{}
	err = json.Unmarshal(body, &blockListPage)
	if err != nil {
		return globals.BlockListOutput{}, err
	}
	if len(blockListPage.Data.Blocklist) > 0 {
		blockListOutput.Items = append(blockListOutput.Items, blockListPage.Data.Blocklist...)
	} else {
		return globals.BlockListOutput{}, err
	}

	// Now we cycle through /2, /3, etc until there are no more
	// https://api.clearsky.services/api/v1/anon/single-blocklist/did:plc:ccskhvd467uwdrxpwaudnbni/2
	if len(blockListOutput.Items) >= 100 {
		for i := 2; i <= maxPages; i++ {
			url = fmt.Sprintf("https://api.clearsky.services/api/v1/anon/single-blocklist/%s/%d", userId, i)
			body, err = FetchUrl(url, logger)
			if err != nil {
				return globals.BlockListOutput{}, err
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
	return blockListOutput, nil
}

func ExpandBlockListUsers(blockList *[]globals.BlockingUser, batchOperationTimeout int, logger *logrus.Logger) (err error) {
	var (
		blockingUser globals.BlockingUser
		userObject   globals.BlueSkyUser
		url          string
		requester    *insrequester.Request
		workerCount  = 100
		wg           sync.WaitGroup
	)
	requester = insrequester.NewRequester().Load()
	requester.WithTimeout(time.Duration(batchOperationTimeout) * time.Second)
	jobs := make(chan globals.Job, len(*blockList))
	results := make(chan *http.Response, len(*blockList))

	for w := 0; w < workerCount; w++ {
		go worker(requester, jobs, results, &wg, logger)
	}

	wg.Add(len(*blockList))
	for _, blockingUser = range *blockList {
		url = fmt.Sprintf("https://public.api.bsky.app/xrpc/app.bsky.actor.getProfile?actor=%s", (blockingUser.DID))
		jobs <- globals.Job{URL: url}
	}
	close(jobs)
	wg.Wait()

	for i := 0; i < len(*blockList); i++ {
		resp := <-results
		if resp != nil {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				fmt.Println(string(body))
				os.Exit(0)
			}

			userObject = globals.BlueSkyUser{}
			err = json.Unmarshal(body, &userObject)
			if err != nil {
				panic(err)
			}
			(*blockList)[i].Username = userObject.Handle
			(*blockList)[i].DisplayName = gomoji.RemoveEmojis(userObject.DisplayName)
			(*blockList)[i].Description = userObject.Description
			(*blockList)[i].Banner = userObject.Banner
			(*blockList)[i].FollowsCount = userObject.FollowsCount
			(*blockList)[i].FollowersCount = userObject.FollowersCount
			(*blockList)[i].Posts = userObject.Posts
		} else {
			fmt.Println("Ouch! nil response")
		}
	}
	return nil
}
