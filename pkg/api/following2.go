package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/gdanko/clearsky/globals"
	"github.com/kr/pretty"
	"github.com/sirupsen/logrus"
)

type FollowingValue struct {
	Type      string `json:"$type"`
	Subject   string `json:"subject"`
	CreatedAt string `json:"createdAt"`
}

type Following2 struct {
	URI   string         `json:"uri"`
	CID   string         `json:"cid"`
	Value FollowingValue `json:"value"`
}

type FollowingListPage2 struct {
	Error   string       `json:"error"`
	Message string       `json:"message"`
	Cursor  string       `json:"cursor"`
	Records []Following2 `json:"records"`
}

func GetFollowing2(targetCredentials globals.BlueSkyCredentials, showFollowingUsers bool, batchOperationTimeout int, listMaxResults int, logger *logrus.Logger) (followingList map[string]Following2, err error) {
	// https://public.api.bsky.app/xrpc/app.bsky.actor.getProfile?actor=did:plc:pk5ffzlhjuckl5r65pcrryjv
	var (
		body []byte
		// credentials         globals.BlueSkyCredentials
		followingListAll    = map[string]Following2{}
		followingListPage   FollowingListPage2
		followingListUser   Following2
		followingPageCursor string
		headers             map[string]string
		listRecordsLimit    = 30
		url                 string
	)
	// credentials = globals.GetCredentials()
	followingListAll = make(map[string]Following2)

	for {
		url = fmt.Sprintf("%s/xrpc/com.atproto.repo.listRecords?repo=%s&collection=app.bsky.graph.follow&limit=%d&cursor=%s", targetCredentials.ServiceEndpoint, targetCredentials.UserDid, listRecordsLimit, followingPageCursor)
		body, err = FetchUrl("GET", url, headers, nil, logger)
		if err != nil {
			return followingList, err
		}
		followingListPage = FollowingListPage2{}
		err = json.Unmarshal(body, &followingListPage)
		if err != nil {
			return followingList, err
		}

		if len(followingListPage.Error) > 0 && len(followingListPage.Message) > 0 {
			return followingList, errors.New(fmt.Sprintf("Could not retrieve a list of followed users: %s", followingListPage.Message))
		}

		// fmt.Println(len(followingListPage.Records))
		if len(followingListPage.Records) > 0 {
			for _, followingListUser = range followingListPage.Records {
				followingListAll[followingListUser.Value.Subject] = followingListUser
			}
			if followingListPage.Cursor == "" {
				break
			} else {
				followingPageCursor = followingListPage.Cursor
			}
		} else {
			break
		}
	}
	pretty.Println(followingListAll)
	fmt.Println(len(followingListAll))
	os.Exit(0)
	return followingListAll, nil
}
