package main

import (
	"fmt"
	"time"

	"github.com/lukeelten/micc-creation-day-26/backend/pkg/client"
	"github.com/lukeelten/micc-creation-day-26/backend/pkg/models"
	"github.com/lukeelten/micc-creation-day-26/backend/pkg/utils"
)

type Task func(client *client.Client, runId string, taskTime time.Duration) error

var simuTasks map[string]Task

func init() {
	simuTasks = make(map[string]Task, 4)
	simuTasks[utils.TASK_DOWNLOAD] = downloadRunDataTask
	simuTasks[utils.TASK_CONVERT] = convertFormatsTask
	simuTasks[utils.TASK_PROCESS] = processDataTask
	simuTasks[utils.TASK_UPLOAD] = uploadResultsTask
}

func downloadRunDataTask(c *client.Client, runId string, taskTime time.Duration) error {
	c.UpdateRunStatusProcessing(runId)
	return genericRunTask("Download Run Data", c, runId, taskTime)
}

func convertFormatsTask(c *client.Client, runId string, taskTime time.Duration) error {
	return genericRunTask("Convert Data Format", c, runId, taskTime)
}

func processDataTask(c *client.Client, runId string, taskTime time.Duration) error {
	return genericRunTask("Process Data", c, runId, taskTime)
}

func uploadResultsTask(c *client.Client, runId string, taskTime time.Duration) error {
	return genericRunTask("Uploading Results", c, runId, taskTime)
}

func genericRunTask(taskName string, c *client.Client, runId string, taskTime time.Duration) error {
	c.Logger.Info("Starting Task", "taskName", taskName, "runId", runId, "duration", taskTime.String())

	event := &models.EventsRecord{
		Type:  models.EventsTypeInfo,
		Title: fmt.Sprintf("Task %s started", taskName),
	}
	c.Logger.Debug("Logging event", "event", event)
	err := c.CreateEvent(runId, event)
	if err != nil {
		return err
	}

	utils.SimulateWork(c.Ctx, taskTime)

	event = &models.EventsRecord{
		Type:  models.EventsTypeInfo,
		Title: fmt.Sprintf("Task %s completed", taskName),
	}

	c.Logger.Debug("Logging event", "event", event)
	c.Logger.Info("Completed Task", "taskName", taskName, "runId", runId, "duration", taskTime.String())
	return c.CreateEvent(runId, event)
}
