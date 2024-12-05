package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gdanko/clearsky/globals"
	"github.com/sirupsen/logrus"
)

type FollowerViewer struct {
	BlockedBy  bool   `json:"blockedBy"`
	FollowedBy string `json:"followedBy"`
	Muted      bool   `json:"muted"`
}

type Follower struct {
	Avatar      string         `json:"avatar"`
	CreatedAt   string         `json:"createdAt"`
	Description string         `json:"description"`
	DID         string         `json:"did"`
	DisplayName string         `json:"displayName"`
	Handle      string         `json:"handle"`
	IndexedAt   string         `json:"indexedAt"`
	Viewer      FollowerViewer `json:"viewer"`
}

type FollowersListSubject struct {
	Avatar      string `json:"avatar"`
	CreatedAt   string `json:"createdAt"`
	DID         string `json:"did"`
	DisplayName string `json:"displayName"`
	Handle      string `json:"handle"`
	IndexedAt   string `json:"indexedAt"`
}

type FollowersListPage struct {
	Error     string               `json:"error"`
	Message   string               `json:"message"`
	Cursor    string               `json:"cursor"`
	Followers []Follower           `json:"followers"`
	Subject   FollowersListSubject `json:"subject"`
}

func GetFollowers(targetCredentials globals.BlueSkyCredentials, showFollowingUsers bool, batchOperationTimeout int, listMaxResults int, logger *logrus.Logger) (followersList map[string]Follower, err error) {
	// https://public.api.bsky.app/xrpc/app.bsky.actor.getProfile?actor=did:plc:pk5ffzlhjuckl5r65pcrryjv
	var (
		body                []byte
		credentials         globals.BlueSkyCredentials
		followersListAll    = map[string]Follower{}
		followersListPage   FollowersListPage
		followersListUser   Follower
		followersPageCursor string
		headers             map[string]string
		listRecordsLimit    = 100
		url                 string
	)
	credentials = globals.GetCredentials()
	followersListAll = make(map[string]Follower)

	for {
		// https://{targetCredentials.ServiceEndpoint}/xrpc/com.atproto.repo.listRecords?repo={targetCredentials.UserDid}&limit=100&collection=app.bsky.graph.follow&cursor=
		// This is anonymous and returns the number of followers reported on the profile page but it does not return the handle and thus requires a did lookup of every entry
		// However, some of these users may be deleted

		// "https://{credentials.ServiceEndpoint}/xrpc/app.bsky.graph.getFollowers?actor={targetCredentials.UserDid}&limit=100&cursor="
		// This requires authentication and returns fewer followers than reported on the profile page. It returns the handle in the results. It probably excludes deleted users.
		url = fmt.Sprintf("%s/xrpc/app.bsky.graph.getFollowers?actor=%s&limit=%d&cursor=%s", credentials.ServiceEndpoint, targetCredentials.UserDid, listRecordsLimit, followersPageCursor)
		headers = map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", credentials.AccessJwtCookie),
			"Content-Type":  "appliation/json",
		}

		body, err = FetchUrl("GET", url, headers, nil, logger)
		if err != nil {
			return followersList, err
		}
		followersListPage = FollowersListPage{}
		err = json.Unmarshal(body, &followersListPage)
		if err != nil {
			return followersList, err
		}

		if len(followersListPage.Error) > 0 && len(followersListPage.Message) > 0 {
			return followersList, errors.New(fmt.Sprintf("Could not retrieve a list of followers: %s", followersListPage.Message))
		}

		if len(followersListPage.Followers) > 0 {
			for _, followersListUser = range followersListPage.Followers {
				followersListAll[followersListUser.DID] = followersListUser
			}
			if followersListPage.Cursor == "" {
				break
			} else {
				followersPageCursor = followersListPage.Cursor
			}
		} else {
			break
		}
	}
	return followersListAll, nil
}
