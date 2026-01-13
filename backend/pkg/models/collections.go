package models

const (
	CollectionsAuthorigins   = "_authOrigins"
	CollectionsExternalauths = "_externalAuths"
	CollectionsMfas          = "_mfas"
	CollectionsOtps          = "_otps"
	CollectionsSuperusers    = "_superusers"
	CollectionsEvents        = "events"
	CollectionsRuns          = "runs"
	CollectionsStates        = "states"
	CollectionsUsers         = "users"
)

type UsersRecord struct {
	Avatar          string `json:"avatar,omitempty"`
	Created         string `json:"created"`
	Email           string `json:"email"`
	EmailVisibility bool   `json:"emailVisibility,omitempty"`
	ID              string `json:"id"`
	Name            string `json:"name,omitempty"`
	Password        string `json:"password"`
	TokenKey        string `json:"tokenKey"`
	Updated         string `json:"updated"`
	Verified        bool   `json:"verified,omitempty"`
}
