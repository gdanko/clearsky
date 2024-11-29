package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gdanko/clearsky/globals"
	"github.com/gdanko/clearsky/util"
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
		logger.Debugf("Fetching %s", job.URL)
		res, err := requester.Get(insrequester.RequestEntity{Endpoint: job.URL})
		if err != nil {
			fmt.Println(err)
		}
		results <- res
		wg.Done()
	}
}

func ExpandBlockListUsers(ids []string, batchOperationTimeout int, logger *logrus.Logger) (blockListUsers []globals.BlueSkyUser, err error) {
	// https://medium.com/insiderengineering/concurrent-http-requests-in-golang-best-practices-and-techniques-f667e5a19dea
	var (
		id          string
		userObject  globals.BlueSkyUser
		url         string
		requester   *insrequester.Request
		workerCount = 100
		wg          sync.WaitGroup
	)
	requester = insrequester.NewRequester().Load()
	requester.WithTimeout(time.Duration(batchOperationTimeout) * time.Second)
	jobs := make(chan globals.Job, len(ids))
	results := make(chan *http.Response, len(ids))

	for w := 0; w < workerCount; w++ {
		go worker(requester, jobs, results, &wg, logger)
	}

	wg.Add(len(ids))
	for _, id = range ids {
		url = fmt.Sprintf("https://public.api.bsky.app/xrpc/app.bsky.actor.getProfile?actor=%s", id)
		jobs <- globals.Job{URL: url}
	}
	close(jobs)
	wg.Wait()

	for i := 0; i < len(ids); i++ {
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
			blockListUsers = append(blockListUsers, userObject)
		}
	}
	return blockListUsers, err
}

func ProcessBlockingUsersList(blockingList *map[string]globals.BlockingUser, blockListPage globals.BlockListPage, batchChunkSize int, batchOperationTimeout int, logger *logrus.Logger) (err error) {
	var (
		chunk   []string
		didSet  []string
		divided [][]string
		i       int
		temp    globals.BlockingUser
	)
	for _, blockingUser := range blockListPage.Data.Blocklist {
		(*blockingList)[blockingUser.DID] = globals.BlockingUser{DID: blockingUser.DID, BlockedDate: blockingUser.BlockedDate}
		didSet = append(didSet, blockingUser.DID)
	}
	divided = util.SliceChunker(didSet, batchChunkSize)
	for i, chunk = range divided {
		logger.Debugf("Processing chunk %d of %d", i+1, len(divided))
		blockListUsers, err := ExpandBlockListUsers(chunk, batchOperationTimeout, logger)
		if err != nil {
			return err
		}
		for _, blockListUser := range blockListUsers {
			temp = (*blockingList)[blockListUser.DID]

			temp.Banner = blockListUser.Banner
			temp.Description = blockListUser.Description
			temp.DisplayName = blockListUser.DisplayName
			temp.Error = blockListUser.Error
			temp.FollowersCount = blockListUser.FollowersCount
			temp.FollowsCount = blockListUser.FollowsCount
			temp.Labels = blockListUser.Labels
			temp.Message = blockListUser.Message
			temp.PinnedPost = blockListUser.PinnedPost
			temp.Posts = blockListUser.Posts
			temp.Username = blockListUser.Handle

			(*blockingList)[blockListUser.DID] = temp
		}
	}
	return nil
}

func GetBlockingUsersList(userId string, batchChunkSize int, batchOperationTimeout int, logger *logrus.Logger) (blockingList map[string]globals.BlockingUser, err error) {
	// Get the user info BEFORE inserting into the database...
	// When you have each batch of 100, go get their details
	var (
		blockListPage globals.BlockListPage
		body          []byte
		maxPages      = 3
		url           string
	)

	blockingList = make(map[string]globals.BlockingUser)

	// https://api.clearsky.services/api/v1/anon/single-blocklist/did:plc:ccskhvd467uwdrxpwaudnbni
	url = fmt.Sprintf("https://api.clearsky.services/api/v1/anon/single-blocklist/%s", userId)
	body, err = FetchUrl(url, logger)
	if err != nil {
		return blockingList, err
	}
	blockListPage = globals.BlockListPage{}
	err = json.Unmarshal(body, &blockListPage)
	if err != nil {
		return blockingList, err
	}
	err = ProcessBlockingUsersList(&blockingList, blockListPage, batchChunkSize, batchOperationTimeout, logger)
	if err != nil {
		return blockingList, err
	}

	if len(blockingList) >= 100 {
		for i := 2; i <= maxPages; i++ {
			url = fmt.Sprintf("https://api.clearsky.services/api/v1/anon/single-blocklist/%s/%d", userId, i)
			body, err = FetchUrl(url, logger)
			if err != nil {
				return blockingList, err
			}
			blockListPage = globals.BlockListPage{}
			err = json.Unmarshal(body, &blockListPage)
			if err != nil {
				return blockingList, err
			}
			err = ProcessBlockingUsersList(&blockingList, blockListPage, batchChunkSize, batchOperationTimeout, logger)
			if err != nil {
				return blockingList, err
			}
		}
	}
	return blockingList, nil
}
