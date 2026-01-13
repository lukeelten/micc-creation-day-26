package controller

import (
	"log/slog"
	"slices"

	"github.com/lukeelten/micc-creation-day-26/backend/pkg/models"
	"github.com/pocketbase/dbx"
)

type RunInstance struct {
	Logger *slog.Logger
	RC     *RunController
	Record *models.RunsRecord
}

func NewRunInstance(rc *RunController, record *models.RunsRecord) *RunInstance {
	return &RunInstance{
		RC:     rc,
		Logger: rc.Logger,
		Record: record,
	}
}

func (ri *RunInstance) GetAllStates() ([]*models.StatesRecord, error) {
	records, err := ri.RC.Pb.FindRecordsByFilter(models.CollectionsStates, "run = {:run}", "", 0, 0, dbx.Params{
		"run": ri.Record.ID,
	})

	if err != nil {
		return nil, err
	}

	var states []*models.StatesRecord = make([]*models.StatesRecord, 0)
	for _, r := range records {
		state := &models.StatesRecord{}
		err = state.FromRecord(r)
		if err != nil {
			return nil, err
		}

		states = append(states, state)
	}

	return states, nil
}

func (ri *RunInstance) GetActiveState() (*models.StatesRecord, error) {
	records, err := ri.GetAllStates()
	if err != nil {
		return nil, err
	}

	for _, state := range records {
		if state.Completed == nil {
			return state, nil
		}
	}

	return nil, nil
}

func (ri *RunInstance) GetLastCompletedState() (*models.StatesRecord, error) {
	records, err := ri.GetAllStates()
	if err != nil {
		return nil, err
	}

	for _, state := range records {
		if state.Completed == nil {
			return state, nil
		}
	}

	slices.SortFunc(records, func(a, b *models.StatesRecord) int {
		if a.Completed.Before(*b.Completed) {
			return 1
		} else if a.Completed.After(*b.Completed) {
			return -1
		}
		return 0
	})

	if len(records) > 0 {
		return records[len(records)-1], nil
	}

	return nil, nil
}

func (ri *RunInstance) Start() {
	if ri.Record.Status == models.RunsStatusCompleted || ri.Record.Status == models.RunsStatusFailed {
		ri.Logger.Info("RunInstance already completed or failed, skipping", "runID", ri.Record.ID)
		return
	}

	ri.Logger.Info("Starting controller for run", "runID", ri.Record.ID, "status", ri.Record.Status)

}
