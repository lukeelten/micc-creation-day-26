package controller

import (
	"github.com/lukeelten/micc-creation-day-26/backend/pkg/models"
	"github.com/lukeelten/micc-creation-day-26/backend/pkg/utils"
	"github.com/pocketbase/pocketbase/core"
)

func (rc *RunController) ProcessRun(runId string) error {
	// Fetch record
	runRecord, err := rc.Pb.FindRecordById(models.CollectionsRuns, runId)
	if err != nil {
		rc.Logger.Error("Failed to fetch run record", "runId", runId, "error", err)
		return err
	}

	runModel := models.ConvertRunRecord(runRecord)

	switch runModel.Status {
	case models.RunsStatusCreated:
		// Create first task
		err := rc.runCreated(runRecord)
		if err != nil {
			rc.Logger.Error("Failed to create first task for run", "runId", runId, "error", err)
			return err
		}

		return nil

	case models.RunsStatusScheduled:
		// Do intentionally nothing
		return nil

	case models.RunsStatusProcessing:
		// Check for next state
		err := rc.runProcessing(runRecord)
		if err != nil {
			rc.Logger.Error("Failed to process run", "runId", runId, "error", err)
			return err
		}

		return nil

	case models.RunsStatusFailed:
		fallthrough
	case models.RunsStatusCompleted:
		if runModel.RuntimeSeconds > 0 {
			// already processed
			return nil
		}

		allStates, err := rc.GetAllStates(runId)
		if err != nil {
			rc.Logger.Error("Failed to fetch states for completed run", "runId", runId, "error", err)
			return err
		}

		for _, state := range allStates {
			if state.Completed == nil {
				err = rc.SetStateComplete(state.ID)
				if err != nil {
					rc.Logger.Error("Failed to set state complete for completed run", "runId", runId, "stateId", state.ID, "error", err)
					return err
				}
			}
		}

		latestState, err := rc.GetLastCompletedState(runId)
		if err != nil {
			rc.Logger.Error("Failed to fetch last completed state for completed run", "runId", runId, "error", err)
			return err
		}

		if latestState == nil || latestState.Completed == nil {
			rc.Logger.Warn("No completed states found for completed run", "runId", runId)
			return nil
		}

		startingTime := runRecord.GetDateTime("created").Time()
		completionTime := *latestState.Completed
		totalDuration := completionTime.Sub(startingTime)

		runRecord.Set("runtimeSeconds", int(totalDuration.Seconds()))
		err = rc.Pb.Save(runRecord)
		if err != nil {
			rc.Logger.Error("Failed to update run runtime for completed run", "runId", runId, "error", err)
			return err
		}

	default:
		rc.Logger.Warn("Unknown run status", "runId", runId, "status", runModel.Status)
		return nil
	}
	return nil
}

func (rc *RunController) runCreated(runRecord *core.Record) error {
	runId := runRecord.Id

	_, err := rc.CreateKubernetesJob(runId, utils.TASK_DOWNLOAD, utils.RandomDurationLimit(utils.TASK_DOWNLOAD_MAX_DURATION))
	if err != nil {
		rc.Logger.Error("Failed to create Kubernetes job for download task", "runId", runId, "error", err)
		return err
	}

	runRecord.Set("status", models.RunsStatusScheduled)
	err = rc.Pb.Save(runRecord)
	if err != nil {
		rc.Logger.Error("Failed to update run status to scheduled", "runId", runId, "error", err)
		return err
	}

	rc.Logger.Debug("Scheduled run", "runId", runId)

	return nil
}

func (rc *RunController) runProcessing(runRecord *core.Record) error {
	return nil
}
