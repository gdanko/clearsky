package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

// https://medium.com/insiderengineering/concurrent-http-requests-in-golang-best-practices-and-techniques-f667e5a19dea
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

func processBlockingUsersList(blockingList *map[string]globals.BlockingUser, batchOperationTimeout int, logger *logrus.Logger) (err error) {
	var (
		batchChunkSize = 25
		blockingUser   globals.BlockingUser
		blueSkyUser    globals.BlueSkyUser
		chunk          []string
		did            string
		didList        []string
		divided        = [][]string{}
		i              int
		requester      *insrequester.Request
		url            string
		urls           = []string{}
		usersList      globals.BlueSkyUsers
		wg             sync.WaitGroup
		workerCount    = 100
	)

	for _, blockingUser := range *blockingList {
		didList = append(didList, blockingUser.DID)
	}

	divided = util.SliceChunker(didList, batchChunkSize)
	for _, chunk = range divided {
		for i, did = range chunk {
			chunk[i] = fmt.Sprintf("actors=%s", did)
		}
		urls = append(
			urls,
			fmt.Sprintf(
				"https://public.api.bsky.app/xrpc/app.bsky.actor.getProfiles?%s",
				strings.Join(chunk, "&"),
			),
		)
	}

	requester = insrequester.NewRequester().Load()
	requester.WithTimeout(time.Duration(batchOperationTimeout) * time.Second)
	jobs := make(chan globals.Job, len(urls))
	results := make(chan *http.Response, len(urls))

	for w := 0; w < workerCount; w++ {
		go worker(requester, jobs, results, &wg, logger)
	}

	wg.Add(len(urls))
	for _, url = range urls {
		jobs <- globals.Job{URL: url}
	}
	close(jobs)
	wg.Wait()

	for i := 0; i < len(urls); i++ {
		resp := <-results
		if resp != nil {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				fmt.Println(string(body))
				os.Exit(0)
			}
			usersList = globals.BlueSkyUsers{}
			err = json.Unmarshal(body, &usersList)
			if err != nil {
				return err
			}

			for _, blueSkyUser = range usersList.Profiles {
				did = blueSkyUser.DID
				blockingUser = (*blockingList)[did]

				blockingUser.Banner = blueSkyUser.Banner
				blockingUser.DisplayName = blueSkyUser.DisplayName
				blockingUser.Error = blueSkyUser.Error
				blockingUser.FollowersCount = blueSkyUser.FollowersCount
				blockingUser.FollowsCount = blueSkyUser.FollowsCount
				blockingUser.Handle = blueSkyUser.Handle
				blockingUser.Labels = blueSkyUser.Labels
				blockingUser.Message = blueSkyUser.Message
				blockingUser.PinnedPost = blueSkyUser.PinnedPost
				blockingUser.Posts = blueSkyUser.Posts

				(*blockingList)[did] = blockingUser
			}
		}
	}
	return nil
}

func GetBlockingUsersList(userId string, batchOperationTimeout int, listMaxResults int, logger *logrus.Logger) (blockingList map[string]globals.BlockingUser, err error) {
	var (
		blockingListAll     = map[string]globals.BlockingUser{}
		blockListPage       globals.BlockListPage
		body                []byte
		limitedBlockingList = map[string]globals.BlockingUser{}
		maxPages            = 1000
		totalRecords        int
		url                 string
	)
	blockingList = map[string]globals.BlockingUser{}
	blockingListAll = make(map[string]globals.BlockingUser)

	url = fmt.Sprintf("https://api.clearsky.services/api/v1/anon/single-blocklist/%s", userId)
	body, err = FetchUrl(url, logger)
	if err != nil {
		return blockingList, err
	}
	blockListPage = globals.BlockListPage{}
	err = json.Unmarshal(body, &blockListPage)
	if err != nil {
		return blockingList, nil
	}
	for _, blockingUser := range blockListPage.Data.Blocklist {
		blockingListAll[blockingUser.DID] = blockingUser
	}

	if len(blockListPage.Data.Blocklist) >= 100 {
		for i := 2; i <= maxPages; i++ {
			url = fmt.Sprintf("https://api.clearsky.services/api/v1/anon/single-blocklist/%s/%d", userId, i)
			body, err = FetchUrl(url, logger)
			if err != nil {
				return blockingList, err
			}
			blockListPage = globals.BlockListPage{}
			err = json.Unmarshal(body, &blockListPage)
			if err != nil {
				return blockingList, nil
			}
			if len(blockListPage.Data.Blocklist) > 0 {
				for _, blockingUser := range blockListPage.Data.Blocklist {
					blockingListAll[blockingUser.DID] = blockingUser
				}
			} else {
				break
			}
		}
	}

	totalRecords = len(blockingListAll)
	if listMaxResults < totalRecords {
		logger.Debugf("Limiting the number of records to %d because the --limit flag was used", listMaxResults)
		limitedBlockingList = make(map[string]globals.BlockingUser)
		for key, value := range blockingListAll {
			limitedBlockingList[key] = value
			if len(limitedBlockingList) == listMaxResults {
				blockingList = limitedBlockingList
				break
			}
		}
	} else {
		blockingList = blockingListAll
	}

	err = processBlockingUsersList(&blockingList, batchOperationTimeout, logger)
	if err != nil {
		return blockingList, err
	}
	return blockingList, nil
}
