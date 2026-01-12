package models

import "github.com/pocketbase/pocketbase/core"

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

func ConvertRunRecord(record *core.Record) *RunsRecord {
	return &RunsRecord{
		Author:         record.GetString("author"),
		Created:        record.GetString("created"),
		ID:             record.Id,
		Message:        record.GetString("message"),
		RuntimeSeconds: record.GetInt("runtimeSeconds"),
		Status:         RunsStatusOptions(record.GetString("status")),
		Updated:        record.GetString("updated"),
	}
}
