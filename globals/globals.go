package globals

import "sync"

// User count data
type CountBlock struct {
	DisplayName string `json:"displayName"`
	Value       string `json:"value"`
}

type DataBlock struct {
	AsOf         string     `json:"as of"`
	ActiveCount  CountBlock `json:"active_count"`
	DeletedCount CountBlock `json:"deleted_count"`
	TotalCount   CountBlock `json:"total_count"`
}

type CountData struct {
	Data DataBlock `json:"data"`
}

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

// ClearSky user block
type ClearSkyUserObject struct {
	AvatarUrl        string `json:"avatar_url"`
	HandleIdentifier string `json:"handle_identifier"`
	Handle           string `json:"identifier"`
	PDS              string `json:"pds"`
	UserUrl          string `json:"user_url"`
}

type ClearSkyUser struct {
	Data ClearSkyUserObject `json:"data"`
}

// Concurrent worker job
type Job struct {
	URL string
}

var (
	debugFlag bool
	mu        sync.RWMutex
)

func SetDebugFlag(x bool) {
	mu.Lock()
	debugFlag = x
	mu.Unlock()
}

func GetDebugFlag() (x bool) {
	mu.Lock()
	x = debugFlag
	mu.Unlock()
	return x
}
