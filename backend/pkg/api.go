package pkg

import (
	"net/http"

	"github.com/lukeelten/micc-creation-day-26/backend/pkg/models"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func SetupApi(pb *pocketbase.PocketBase) error {

	pb.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// Setup API routes here
		// Update run status endpoints
		e.Router.GET("/runs/{runId}/complete", func(c *core.RequestEvent) error {
			return updateRunStatus(c, models.RunsStatusCompleted)
		})

		e.Router.GET("/runs/{runId}/failed", func(c *core.RequestEvent) error {
			return updateRunStatus(c, models.RunsStatusFailed)
		})

		e.Router.GET("/runs/{runId}/processing", func(c *core.RequestEvent) error {
			return updateRunStatus(c, models.RunsStatusProcessing)
		})

		e.Router.GET("/runs/{runId}/scheduled", func(c *core.RequestEvent) error {
			return updateRunStatus(c, models.RunsStatusScheduled)
		})

		// Add event endpoint
		e.Router.POST("/runs/{runId}/events", func(c *core.RequestEvent) error {
			return createEvent(c)
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

	record, err := c.App.FindRecordById(models.CollectionsRuns, runId)
	if err != nil {
		return apis.NewNotFoundError("Run not found", err)
	}

	record.Set("status", string(status))

	if err := c.App.Save(record); err != nil {
		return apis.NewApiError(http.StatusInternalServerError, "Failed to update run status", err)
	}

	return c.JSON(http.StatusOK, models.ConvertRunRecord(record))
}

func createEvent(c *core.RequestEvent) error {
	runId := c.Request.PathValue("runId")
	if runId == "" {
		return apis.NewBadRequestError("runId is required", nil)
	}

	// Verify the run exists
	_, err := c.App.FindRecordById(models.CollectionsRuns, runId)
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
	if err := models.CreateEvent(c.App, &eventData); err != nil {
		return apis.NewApiError(http.StatusInternalServerError, "Failed to create event", err)
	}

	return c.JSON(http.StatusCreated, eventData)
}
