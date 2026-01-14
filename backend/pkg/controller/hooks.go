package controller

import (
	"github.com/lukeelten/micc-creation-day-26/backend/pkg/models"
	"github.com/pocketbase/pocketbase/core"
)

func (rc *RunController) SetupHooks() {
	rc.Logger.Debug("Register hooks for run controller")
	pb := rc.Pb

	pb.OnRecordAfterCreateSuccess(models.CollectionsRuns).BindFunc(func(e *core.RecordEvent) error {
		runId := e.Record.Id
		rc.Logger.Info("New run created, adding to work queue", "runId", runId)

		// Add the new run to the work queue
		rc.workQueue <- runId

		return e.Next()
	})

	pb.OnRecordAfterUpdateSuccess(models.CollectionsRuns).BindFunc(func(e *core.RecordEvent) error {
		runId := e.Record.Id
		rc.Logger.Info("Run updated, adding to work queue", "runId", runId)

		// Add the new run to the work queue
		rc.workQueue <- runId

		return e.Next()
	})

	pb.OnRecordAfterCreateSuccess(models.CollectionsStates).BindFunc(func(e *core.RecordEvent) error {
		runId := e.Record.GetString("run")
		rc.Logger.Info("New state created for run", "runId", runId)

		rc.workQueue <- runId

		return e.Next()
	})

	pb.OnRecordAfterUpdateSuccess(models.CollectionsStates).BindFunc(func(e *core.RecordEvent) error {
		runId := e.Record.GetString("run")
		rc.Logger.Info("State updated for run", "runId", runId)

		rc.workQueue <- runId

		return e.Next()
	})
}
