package globals

import (
	"sync"
)

// User's did
type UserData struct {
	AvatarUrl     string `json:"avatar_url"`
	DidIdentifier string `json:"did_identifier"`
	Identifier    string `json:"identifier"`
	PDS           string `json:"pds"`
	UserUrl       string `json:"user_url"`
}

type UserDid struct {
	Data UserData `json:"data"`
}

// Block list structs
type BlockingUser struct {
	BlockedDate    string `json:"blocked_date"`
	DID            string `json:"did"`
	Status         string `json:"status"`
	Username       string `json:"username"`
	DisplayName    string `json:"displayName"`
	Description    string `json:"description"`
	Banner         string `json:"banner"`
	FollowsCount   int    `json:"followsCount"`
	FollowersCount int    `json:"followersCount"`
	Posts          int    `json:"postsCount"`
}

type BlockingUsers struct {
	Blocklist []BlockingUser `json:"blocklist"`
	ItemCount int
	Pages     int `json:"pages"`
}

type BlockListPage struct {
	Data     BlockingUsers `json:"data"`
	Identity string        `json:"identity"`
	Status   bool          `json:"status"`
}

type BlockListOutput struct {
	Items []BlockingUser `json:"items"`
	Count int            `json:"count"`
}

// BlueSky user block
type BlueSkyUser struct {
	DID            string `json:"did"`
	Handle         string `json:"handle"`
	DisplayName    string `json:"displayName"`
	Description    string `json:"description"`
	Banner         string `json:"banner"`
	FollowsCount   int    `json:"followsCount"`
	FollowersCount int    `json:"followersCount"`
	Posts          int    `json:"postsCount"`
}

var (
	mu              sync.RWMutex
	blockListOutput BlockListOutput
	blueSkyUser     BlueSkyUser
	userDid         UserDid
)

// set and get pairs
func SetBlockListOutput(x BlockListOutput) {
	mu.Lock()
	blockListOutput = x
	mu.Unlock()
}

func GetBlockListOutput() (x BlockListOutput) {
	mu.Lock()
	x = blockListOutput
	mu.Unlock()
	return x
}

func SetBlueSkyUser(x BlueSkyUser) {
	mu.Lock()
	blueSkyUser = x
	mu.Unlock()
}

func GetBlueSkyUser() (x BlueSkyUser) {
	mu.Lock()
	x = blueSkyUser
	mu.Unlock()
	return x
}

func SetUserData(x UserDid) {
	mu.Lock()
	userDid = x
	mu.Unlock()
}

func GetUserData() (x UserDid) {
	mu.Lock()
	x = userDid
	mu.Unlock()
	return x
}
