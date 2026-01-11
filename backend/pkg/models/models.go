package models

type EventsTypeOptions string

const (
	EventsTypeInfo  EventsTypeOptions = "info"
	EventsTypeWarn  EventsTypeOptions = "warn"
	EventsTypeError EventsTypeOptions = "error"
)

type EventsRecord struct {
	AdditionalData interface{}       `json:"additionalData,omitempty"`
	AdditionalText string            `json:"additionalText,omitempty"`
	Created        string            `json:"created"`
	Description    string            `json:"description,omitempty"`
	ID             string            `json:"id"`
	Run            string            `json:"run"`
	Title          string            `json:"title"`
	Type           EventsTypeOptions `json:"type"`
	Updated        string            `json:"updated"`
}

type RunsStatusOptions string

const (
	RunsStatusCreated    RunsStatusOptions = "CREATED"
	RunsStatusFailed     RunsStatusOptions = "FAILED"
	RunsStatusScheduled  RunsStatusOptions = "SCHEDULED"
	RunsStatusCompleted  RunsStatusOptions = "COMPLETED"
	RunsStatusProcessing RunsStatusOptions = "PROCESSING"
)

type RunsRecord struct {
	Author         string            `json:"author,omitempty"`
	Created        string            `json:"created"`
	ID             string            `json:"id"`
	Message        string            `json:"message,omitempty"`
	RuntimeSeconds int               `json:"runtimeSeconds,omitempty"`
	Status         RunsStatusOptions `json:"status,omitempty"`
	Updated        string            `json:"updated"`
}

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
