package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gdanko/clearsky/globals"
	"github.com/sirupsen/logrus"
)

type FollowingViewer struct {
	BlockedBy  bool   `json:"blockedBy"`
	FollowedBy string `json:"followedBy"`
	Following  string `json:"following"`
	Muted      bool   `json:"muted"`
}

type Following struct {
	Avatar      string         `json:"avatar"`
	CreatedAt   string         `json:"createdAt"`
	Description string         `json:"description"`
	DID         string         `json:"did"`
	DisplayName string         `json:"displayName"`
	Handle      string         `json:"handle"`
	IndexedAt   string         `json:"indexedAt"`
	Viewer      FollowerViewer `json:"viewer"`
}

type FollowingListSubject struct {
	Avatar      string `json:"avatar"`
	CreatedAt   string `json:"createdAt"`
	DID         string `json:"did"`
	DisplayName string `json:"displayName"`
	Handle      string `json:"handle"`
	IndexedAt   string `json:"indexedAt"`
}

type FollowingListPage struct {
	Error     string               `json:"error"`
	Message   string               `json:"message"`
	Cursor    string               `json:"cursor"`
	Following []Following          `json:"follows"`
	Subject   FollowingListSubject `json:"subject"`
}

func GetFollowing(targetCredentials globals.BlueSkyCredentials, showFollowingUsers bool, batchOperationTimeout int, listMaxResults int, logger *logrus.Logger) (followingList map[string]Following, err error) {
	// https://public.api.bsky.app/xrpc/app.bsky.actor.getProfile?actor=did:plc:pk5ffzlhjuckl5r65pcrryjv
	var (
		body                []byte
		credentials         globals.BlueSkyCredentials
		followingListAll    = map[string]Following{}
		followingListPage   FollowingListPage
		followingListUser   Following
		followingPageCursor string
		headers             map[string]string
		listRecordsLimit    = 30
		url                 string
	)
	credentials = globals.GetCredentials()
	followingListAll = make(map[string]Following)

	for {
		url = fmt.Sprintf("%s/xrpc/app.bsky.graph.getFollows?actor=%s&limit=%d&cursor=%s", credentials.ServiceEndpoint, targetCredentials.UserDid, listRecordsLimit, followingPageCursor)
		headers = map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", credentials.AccessJwtCookie),
			"Content-Type":  "appliation/json",
		}

		body, err = FetchUrl("GET", url, headers, nil, logger)
		if err != nil {
			return followingList, err
		}
		followingListPage = FollowingListPage{}
		err = json.Unmarshal(body, &followingListPage)
		if err != nil {
			return followingList, err
		}

		if len(followingListPage.Error) > 0 && len(followingListPage.Message) > 0 {
			return followingList, errors.New(fmt.Sprintf("Could not retrieve a list of followed users: %s", followingListPage.Message))
		}

		fmt.Println(len(followingListPage.Following))
		if len(followingListPage.Following) > 0 {
			for _, followingListUser = range followingListPage.Following {
				followingListAll[followingListUser.DID] = followingListUser
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
	return followingListAll, nil
}
