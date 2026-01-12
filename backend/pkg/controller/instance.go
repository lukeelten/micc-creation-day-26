package controller

import (
	"log/slog"

	"github.com/lukeelten/micc-creation-day-26/backend/pkg/models"
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

func (ri *RunInstance) Start() {
	if ri.Record.Status == models.RunsStatusCompleted || ri.Record.Status == models.RunsStatusFailed {
		ri.Logger.Info("RunInstance already completed or failed, skipping", "runID", ri.Record.ID)
		return
	}

	ri.Logger.Info("Starting controller for run", "runID", ri.Record.ID)

}
