package globals

import "sync"

// Concurrent worker job
type Job struct {
	URL string
}

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
	Error   string   `json:"error"`
	Message string   `json:"message"`
	Data    UserData `json:"data"`
}

type BlueSkyCredentials struct {
	AccessJwtCookie string
	ServiceEndpoint string
	UserDid         string
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
	Handle         string             `json:"username"`
	Labels         []BlueSkyUserLabel `json:"labels"`
	Message        string             `json:"message"`
	PinnedPost     BlueSkyPinnedPost  `json:"pinnedPost"`
	Posts          int                `json:"postsCount"`
	Status         string             `json:"status"`
}

type BlockedByUsers struct {
	Blocklist []BlockingUser `json:"blocklist"`
	Pages     int            `json:"pages"`
	ItemCount int
}

type BlockedByListPage struct {
	Data     BlockedByUsers `json:"data"`
	Identity string         `json:"identity"`
	Status   bool           `json:"status"`
}

type BlockedByListOutput struct {
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
	Posts          int                `json:"poptsCount"`
}

// List of users via https://public.api.bsky.app/xrpc/app.bsky.actor.getProfiles
type BlueSkyUsers struct {
	Profiles []BlueSkyUser `json:"profiles"`
}

// PLC directory block
type PlcService struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}

type PlcVerificationMethod struct {
	ID                 string `json:"id"`
	Type               string `json:"type"`
	Controller         string `json:"controller"`
	PublicKeyMultibase string `json:"publicKeyMultibase"`
}

type PlcDirectoryEntry struct {
	Context            []string                `json:"@context"`
	ID                 string                  `json:"id"`
	AlsoKnownAs        []string                `json:"alsoKnownAs"`
	VerificationMethod []PlcVerificationMethod `json:"verificationMethod"`
	Service            []PlcService            `json:"service"`
}

// Blocking list structs
type BlockingListValue struct {
	Type      string `json:"$type"`
	Subject   string `json:"subject"`
	CreatedAd string `json:"createdAt"`
}

type BlockingListRecord struct {
	URI   string            `json:"uri"`
	CID   string            `json:"cid"`
	Value BlockingListValue `json:"value"`
}

type BlockingListPage struct {
	Records []BlockingListRecord `json:"records"`
	Cursor  string               `json:"cursor"`
}

// bsky.app session document
type SessionDocument struct {
	Error           string            `json:"error"`
	Message         string            `json:"message"`
	DID             string            `json:"did"`
	DIDDoc          PlcDirectoryEntry `json:"didDoc"`
	Handle          string            `json:"handle"`
	Email           string            `json:"email"`
	EmailConfirmed  bool              `json:"emailConfirmed"`
	EmailAuthFactor bool              `json:"emailAuthFactor"`
	AccessJwt       string            `json:"accessJwt"`
	RefreshJwt      string            `json:"refreshJwt"`
	Active          bool              `json:"active"`
}

// bsky.app credentials
type Credentials struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

var (
	blueSkyCredentials BlueSkyCredentials
	mu                 sync.RWMutex
)

func SetCredentials(x BlueSkyCredentials) {
	mu.Lock()
	blueSkyCredentials = x
	mu.Unlock()
}

func GetCredentials() (x BlueSkyCredentials) {
	mu.Lock()
	x = blueSkyCredentials
	mu.Unlock()
	return x
}
