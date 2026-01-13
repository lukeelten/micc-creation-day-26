package models

import (
	"time"

	"github.com/pocketbase/pocketbase/core"
)

type StatesTaskOptions string

const (
	StatesTaskDownload StatesTaskOptions = "download"
	StatesTaskConvert  StatesTaskOptions = "convert"
	StatesTaskProcess  StatesTaskOptions = "process"
	StatesTaskUpload   StatesTaskOptions = "upload"
)

type StatesRecord struct {
	ID        string            `json:"id"`
	Created   time.Time         `json:"created"`
	Updated   time.Time         `json:"updated"`
	Completed *time.Time        `json:"completed,omitempty"`
	Run       string            `json:"run,omitempty"`
	Task      StatesTaskOptions `json:"task,omitempty"`
}

func (st *StatesRecord) FromRecord(record *core.Record) error {
	st.ID = record.Id
	st.Created = record.GetDateTime("created").Time()
	st.Updated = record.GetDateTime("updated").Time()
	st.Run = record.GetString("run")
	st.Task = StatesTaskOptions(record.GetString("task"))

	if c := record.GetString("completed"); len(c) > 0 {
		t := record.GetDateTime("completed").Time()
		st.Completed = &t
	}

	return nil
}
