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
	"github.com/useinsider/go-pkg/insrequester"
	// "golang.org/x/sync/errgroup"
)

func GetUserID(accountName string) (userId string, err error) {
	var (
		body []byte
		url  string
	)
	url = fmt.Sprintf("https://api.clearsky.services/api/v1/anon/get-did/%s", accountName)
	body, err = FetchUrl(url)
	if err != nil {
		return userId, err
	}

	getDid := globals.UserDid{}
	err = json.Unmarshal(body, &getDid)
	if err != nil {
		return userId, err
	}
	userId = getDid.Data.DidIdentifier

	return userId, nil
}

func worker(requester *insrequester.Request, jobs <-chan globals.Job, results chan<- *http.Response, wg *sync.WaitGroup) {
	for job := range jobs {
		fmt.Println(job.URL)
		res, err := requester.Get(insrequester.RequestEntity{Endpoint: job.URL})
		if err != nil {
			fmt.Println(err)
		}
		results <- res
		wg.Done()
	}
}

func ExpandBlockListUsers(blockList *[]globals.BlockingUser) (err error) {
	var (
		blockingUser   globals.BlockingUser
		blueSkyUser    globals.BlueSkyUser
		url            string
		requester      *insrequester.Request
		requestTimeout = 30
		workerCount    = 100
		wg             sync.WaitGroup
	)
	requester = insrequester.NewRequester().Load()
	requester.WithTimeout(time.Duration(requestTimeout) * time.Second)
	jobs := make(chan globals.Job, len(*blockList))
	results := make(chan *http.Response, len(*blockList))

	for w := 0; w < workerCount; w++ {
		go worker(requester, jobs, results, &wg)
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

			blueSkyUser = globals.BlueSkyUser{}
			err = json.Unmarshal(body, &blueSkyUser)
			if err != nil {
				fmt.Println(err)
			}
			(*blockList)[i].Username = blueSkyUser.Handle
			(*blockList)[i].DisplayName = blueSkyUser.DisplayName
			(*blockList)[i].Description = blueSkyUser.Description
			(*blockList)[i].Banner = blueSkyUser.Banner
			(*blockList)[i].FollowsCount = blueSkyUser.FollowsCount
			(*blockList)[i].FollowersCount = blueSkyUser.FollowersCount
			(*blockList)[i].Posts = blueSkyUser.Posts
		} else {
			fmt.Println("Ouch! nil response")
		}
	}
	return nil
}
