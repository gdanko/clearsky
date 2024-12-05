package api

import (
	"encoding/json"
	"errors"
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

var (
	body        []byte
	credentials globals.Credentials
	headers     map[string]string
	jsonBytes   []byte
	userDid     globals.UserDid
	url         string
)

// Get the DID from the handle
func GetUserDid(blueSkyHandle string, logger *logrus.Logger) (targetDid string, err error) {
	url = fmt.Sprintf("https://api.clearsky.services/api/v1/anon/get-did/%s", blueSkyHandle)
	body, err = FetchUrl("GET", url, map[string]string{}, nil, logger)
	if err != nil {
		return targetDid, err
	}
	userDid = globals.UserDid{}
	err = json.Unmarshal(body, &userDid)
	if err != nil {
		return targetDid, err
	}
	if userDid.Error != "" {
		return targetDid, errors.New(fmt.Sprintf("Failed to get DID information for %s: %s", blueSkyHandle, userDid.Error))
	}
	return userDid.Data.DidIdentifier, nil
}

// Get the PLC directory info
func GetPlcDirectoryInfo(userId string, logger *logrus.Logger) (plcDirectoryEntry globals.PlcDirectoryEntry, err error) {
	url = fmt.Sprintf("https://plc.directory/%s", userId)
	body, err = FetchUrl("GET", url, map[string]string{}, nil, logger)
	if err != nil {
		return plcDirectoryEntry, err
	}
	plcDirectoryEntry = globals.PlcDirectoryEntry{}
	err = json.Unmarshal(body, &plcDirectoryEntry)
	if err != nil {
		return plcDirectoryEntry, err
	}
	return plcDirectoryEntry, nil
}

// Create the session and get the cookie
func CreateBlueSkySession(blueSkyHandle, blueSkyPassword, serviceEndpoint string, logger *logrus.Logger) (sessionDocument globals.SessionDocument, err error) {
	url = fmt.Sprintf("%s/xrpc/com.atproto.server.createSession", serviceEndpoint)
	credentials = globals.Credentials{Identifier: blueSkyHandle, Password: blueSkyPassword}
	headers = map[string]string{
		"Content-Type": "application/json",
	}
	jsonBytes, err = json.Marshal(credentials)
	if err != nil {
		return sessionDocument, err
	}
	body, err = FetchUrl("POST", url, headers, jsonBytes, logger)
	if err != nil {
		return sessionDocument, err
	}

	sessionDocument = globals.SessionDocument{}
	err = json.Unmarshal(body, &sessionDocument)
	if err != nil {
		return sessionDocument, err
	}
	if sessionDocument.Error != "" && sessionDocument.Message != "" {
		return sessionDocument, errors.New(sessionDocument.Message)
	}
	return sessionDocument, nil
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

func processUsersList(userList *map[string]globals.BlockingUser, batchOperationTimeout int, logger *logrus.Logger) (err error) {
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

	for _, blockingUser := range *userList {
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
				blockingUser = (*userList)[did]

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

				(*userList)[did] = blockingUser
			}
		}
	}
	return nil
}
