package controller

import (
	"slices"
	"time"

	"github.com/lukeelten/micc-creation-day-26/backend/pkg/models"
	"github.com/lukeelten/micc-creation-day-26/backend/pkg/utils"
	"github.com/pocketbase/dbx"
)

func (rc *RunController) GetAllStates(runId string) ([]*models.StatesRecord, error) {
	records, err := rc.Pb.FindRecordsByFilter(models.CollectionsStates, "run = {:run}", "", 0, 0, dbx.Params{
		"run": runId,
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

func (rc *RunController) SetStateComplete(stateId string) error {
	stateRecord, err := rc.Pb.FindRecordById(models.CollectionsStates, stateId)
	if err != nil {
		return err
	}

	if len(stateRecord.GetString("completed")) == 0 {
		stateRecord.Set("completed", time.Now().String())
		err = rc.Pb.Save(stateRecord)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rc *RunController) GetActiveState(runId string) (*models.StatesRecord, error) {
	records, err := rc.GetAllStates(runId)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(records, utils.SortStatesByCompletedDesc)

	for _, state := range records {
		if state.Completed == nil {
			return state, nil
		}
	}

	return nil, nil
}

func (rc *RunController) GetLastCompletedState(runId string) (*models.StatesRecord, error) {
	records, err := rc.GetAllStates(runId)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(records, utils.SortStatesByCompletedDesc)

	if len(records) > 0 {
		return records[0], nil
	}

	return nil, nil
}
