package globals

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
	Banner         string             `json:"banner"`
	BlockedDate    string             `json:"blocked_date"`
	Description    string             `json:"description"`
	DID            string             `json:"did"`
	DisplayName    string             `json:"displayName"`
	Error          string             `json:"error"`
	FollowersCount int                `json:"followersCount"`
	FollowsCount   int                `json:"followsCount"`
	Labels         []BlueSkyUserLabel `json:"labels"`
	Message        string             `json:"message"`
	PinnedPost     BlueSkyPinnedPost  `json:"pinnedPost"`
	Posts          int                `json:"postsCount"`
	Status         string             `json:"status"`
	Username       string             `json:"username"`
}

type BlockingUsers struct {
	Blocklist []BlockingUser `json:"blocklist"`
	Pages     int            `json:"pages"`
	ItemCount int
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

type BlueSkyPinnedPost struct {
	CID string `json:"cid"`
	URI string `json:"uri"`
}

type BlueSkyUserLabel struct {
	SRC string `json:"src"`
	URI string `json:"url"`
	CID string `json:"cid"`
	VAL string `json:"val"`
	CTS string `json:"cts"`
}

// BlueSky user block
type BlueSkyUser struct {
	Banner         string             `json:"banner"`
	Description    string             `json:"description"`
	DID            string             `json:"did"`
	DisplayName    string             `json:"displayName"`
	Error          string             `json:"error"`
	FollowersCount int                `json:"followersCount"`
	FollowsCount   int                `json:"followsCount"`
	Handle         string             `json:"handle"`
	Labels         []BlueSkyUserLabel `json:"labels"`
	Message        string             `json:"message"`
	PinnedPost     BlueSkyPinnedPost  `json:"pinnedPost"`
	Posts          int                `json:"postsCount"`
}

// Concurrent worker job
type Job struct {
	URL string
}
