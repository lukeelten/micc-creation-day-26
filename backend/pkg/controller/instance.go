package controller

import (
	"log/slog"

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

func (ri *RunInstance) Start() {
	if ri.Record.Status == models.RunsStatusCompleted || ri.Record.Status == models.RunsStatusFailed {
		ri.Logger.Info("RunInstance already completed or failed, skipping", "runID", ri.Record.ID)
		return
	}

	ri.Logger.Info("Starting controller for run", "runID", ri.Record.ID, "status", ri.Record.Status)

	// Find states for this run
	states, err := ri.RC

}
