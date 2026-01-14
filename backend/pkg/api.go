package pkg

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/lukeelten/micc-creation-day-26/backend/pkg/models"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func SetupApi(pb *pocketbase.PocketBase) error {
	pb.Logger().Info("Setting up API routes...")

	pb.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// Setup API routes here
		// Update run status endpoints
		e.Router.GET("/v1/runs/{runId}/complete", func(c *core.RequestEvent) error {
			return updateRunStatus(c, models.RunsStatusCompleted)
		})

		e.Router.GET("/v1/runs/{runId}/failed", func(c *core.RequestEvent) error {
			return updateRunStatus(c, models.RunsStatusFailed)
		})

		e.Router.GET("/v1/runs/{runId}/processing", func(c *core.RequestEvent) error {
			return updateRunStatus(c, models.RunsStatusProcessing)
		})

		e.Router.GET("/v1/runs/{runId}/scheduled", func(c *core.RequestEvent) error {
			return updateRunStatus(c, models.RunsStatusScheduled)
		})

		// Add event endpoint
		e.Router.POST("/v1/runs/{runId}/events", func(c *core.RequestEvent) error {
			return createEvent(c)
		})

		// States endpoints
		e.Router.GET("/v1/states/{runId}/{task}/start", func(c *core.RequestEvent) error {
			return startState(c)
		})

		e.Router.GET("/v1/states/{runId}/{task}/stop", func(c *core.RequestEvent) error {
			return stopState(c)
		})

		return e.Next()
	})

	return nil
}

func updateRunStatus(c *core.RequestEvent, status models.RunsStatusOptions) error {
	runId := c.Request.PathValue("runId")
	if runId == "" {
		return apis.NewBadRequestError("runId is required", nil)
	}

	return c.App.RunInTransaction(func(txApp core.App) error {
		record, err := txApp.FindRecordById(models.CollectionsRuns, runId)
		if err != nil {
			return apis.NewNotFoundError("Run not found", err)
		}

		record.Set("status", string(status))

		if err := txApp.Save(record); err != nil {
			return apis.NewApiError(http.StatusInternalServerError, "Failed to update run status", err)
		}

		return c.JSON(http.StatusOK, models.ConvertRunRecord(record))
	})
}

func createEvent(c *core.RequestEvent) error {
	runId := c.Request.PathValue("runId")
	if runId == "" {
		return apis.NewBadRequestError("runId is required", nil)
	}

	return c.App.RunInTransaction(func(txApp core.App) error {
		// Verify the run exists
		_, err := txApp.FindRecordById(models.CollectionsRuns, runId)
		if err != nil {
			return apis.NewNotFoundError("Run not found", err)
		}

		// Parse the event data from request body
		var eventData models.EventsRecord
		if err := c.BindBody(&eventData); err != nil {
			return apis.NewBadRequestError("Invalid request body", err)
		}

		// Set the run ID
		eventData.Run = runId

		// Create the event
		if err := models.CreateEvent(txApp, &eventData); err != nil {
			return apis.NewApiError(http.StatusInternalServerError, "Failed to create event", err)
		}

		return c.JSON(http.StatusCreated, eventData)
	})
}

func startState(c *core.RequestEvent) error {
	runId := c.Request.PathValue("runId")
	task := c.Request.PathValue("task")

	if runId == "" {
		return apis.NewBadRequestError("runId is required", nil)
	}

	if task == "" {
		return apis.NewBadRequestError("task is required", nil)
	}

	c.App.Logger().Debug("/states/runId/task/start", "runId", runId, "task", task)

	return c.App.RunInTransaction(func(txApp core.App) error {
		record, err := txApp.FindFirstRecordByFilter(models.CollectionsStates, "run = {:run} && task = {:task}", dbx.Params{"run": runId, "task": task})
		if err != nil && err != sql.ErrNoRows {
			return apis.NewApiError(http.StatusInternalServerError, "Failed to query run", err)
		}

		if err != sql.ErrNoRows {
			return apis.NewBadRequestError("State already started for this run and task", nil)
		}

		collection, err := txApp.FindCollectionByNameOrId(models.CollectionsStates)
		if err != nil {
			return apis.NewApiError(http.StatusInternalServerError, "Failed to create state record", err)
		}

		record = core.NewRecord(collection)
		record.Set("run", runId)
		record.Set("task", string(task))

		if err := txApp.Save(record); err != nil {
			return apis.NewApiError(http.StatusInternalServerError, "Failed to create state record", err)
		}

		return c.JSON(http.StatusCreated, record)
	})
}

func stopState(c *core.RequestEvent) error {
	runId := c.Request.PathValue("runId")
	task := c.Request.PathValue("task")

	if runId == "" {
		return apis.NewBadRequestError("runId is required", nil)
	}

	if task == "" {
		return apis.NewBadRequestError("task is required", nil)
	}

	c.App.Logger().Debug("/states/runId/task/stop", "runId", runId, "task", task)

	return c.App.RunInTransaction(func(txApp core.App) error {
		record, err := txApp.FindFirstRecordByFilter(models.CollectionsStates, "run = {:run} && task = {:task}", dbx.Params{"run": runId, "task": task})
		if err != nil && err != sql.ErrNoRows {
			return apis.NewApiError(http.StatusInternalServerError, "Failed to query run", err)
		}

		if err == sql.ErrNoRows {
			return apis.NewNotFoundError("State record not found", nil)
		}

		c.App.Logger().Info("Found state to stop", "stateId", record.Id, "runId", runId, "task", task)

		now := time.Now()
		record.Set("completed", now)

		if err := txApp.Save(record); err != nil {
			return apis.NewApiError(http.StatusInternalServerError, "Failed to save state record", err)
		} else {
			c.App.Logger().Info("State stopped successfully", "stateId", record.Id, "runId", runId, "task", task)
		}

		return c.JSON(http.StatusOK, record)
	})
}
