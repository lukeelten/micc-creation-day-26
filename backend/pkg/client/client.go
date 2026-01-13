package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/lukeelten/micc-creation-day-26/backend/pkg/models"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Logger     *slog.Logger
	Ctx        context.Context
}

func NewClient(ctx context.Context, baseURL string, logger *slog.Logger) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
		Logger:     logger,
		Ctx:        ctx,
	}
}

// UpdateRunStatusComplete marks a run as completed
func (c *Client) UpdateRunStatusComplete(runId string) (*models.RunsRecord, error) {
	return c.updateRunStatus(runId, "complete")
}

// UpdateRunStatusFailed marks a run as failed
func (c *Client) UpdateRunStatusFailed(runId string) (*models.RunsRecord, error) {
	return c.updateRunStatus(runId, "failed")
}

// UpdateRunStatusProcessing marks a run as processing
func (c *Client) UpdateRunStatusProcessing(runId string) (*models.RunsRecord, error) {
	return c.updateRunStatus(runId, "processing")
}

// UpdateRunStatusScheduled marks a run as scheduled
func (c *Client) UpdateRunStatusScheduled(runId string) (*models.RunsRecord, error) {
	return c.updateRunStatus(runId, "scheduled")
}

func (c *Client) updateRunStatus(runId, status string) (*models.RunsRecord, error) {
	url := fmt.Sprintf("%s/runs/%s/%s", c.BaseURL, runId, status)

	c.Logger.InfoContext(c.Ctx, "Updating run status", "runId", runId, "status", status, "url", url)

	req, err := http.NewRequestWithContext(c.Ctx, http.MethodGet, url, nil)
	if err != nil {
		c.Logger.ErrorContext(c.Ctx, "Failed to create request", "error", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.Logger.ErrorContext(c.Ctx, "Failed to send request", "error", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.Logger.ErrorContext(c.Ctx, "Unexpected status code", "statusCode", resp.StatusCode, "body", string(body))
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var run models.RunsRecord
	if err := json.NewDecoder(resp.Body).Decode(&run); err != nil {
		c.Logger.ErrorContext(c.Ctx, "Failed to decode response", "error", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	c.Logger.InfoContext(c.Ctx, "Successfully updated run status", "runId", runId, "status", status)
	return &run, nil
}

// CreateEvent creates a new event for a run
func (c *Client) CreateEvent(runId string, event *models.EventsRecord) error {
	event.Run = runId

	url := fmt.Sprintf("%s/runs/%s/events", c.BaseURL, runId)

	c.Logger.InfoContext(c.Ctx, "Creating event", "runId", runId, "eventTitle", event.Title, "eventType", event.Type)

	jsonData, err := json.Marshal(event)
	if err != nil {
		c.Logger.ErrorContext(c.Ctx, "Failed to marshal event", "error", err)
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	req, err := http.NewRequestWithContext(c.Ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.Logger.ErrorContext(c.Ctx, "Failed to create request", "error", err)
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.Logger.ErrorContext(c.Ctx, "Failed to send request", "error", err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		c.Logger.ErrorContext(c.Ctx, "Unexpected status code", "statusCode", resp.StatusCode, "body", string(body))
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	c.Logger.InfoContext(c.Ctx, "Successfully created event", "runId", runId, "eventTitle", event.Title)
	return nil
}

// StartState creates a new state record for a task
func (c *Client) StartState(runId string, task string) (*models.StatesRecord, error) {
	url := fmt.Sprintf("%s/states/%s/%s/start", c.BaseURL, runId, string(task))

	c.Logger.InfoContext(c.Ctx, "Starting state record", "runId", runId, "task", string(task), "url", url)

	req, err := http.NewRequestWithContext(c.Ctx, http.MethodGet, url, nil)
	if err != nil {
		c.Logger.ErrorContext(c.Ctx, "Failed to create request", "error", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.Logger.ErrorContext(c.Ctx, "Failed to send request", "error", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		c.Logger.ErrorContext(c.Ctx, "Unexpected status code", "statusCode", resp.StatusCode, "body", string(body))
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var state models.StatesRecord
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		c.Logger.ErrorContext(c.Ctx, "Failed to decode response", "error", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	c.Logger.InfoContext(c.Ctx, "Successfully started state record", "runId", runId, "task", string(task), "stateId", state.ID)
	return &state, nil
}

// StopState completes a state record
func (c *Client) StopState(runId string, task string) (*models.StatesRecord, error) {
	url := fmt.Sprintf("%s/states/%s/%s/stop", c.BaseURL, runId, string(task))

	c.Logger.InfoContext(c.Ctx, "Stopping state record", "runId", runId, "task", string(task), "url", url)

	req, err := http.NewRequestWithContext(c.Ctx, http.MethodGet, url, nil)
	if err != nil {
		c.Logger.ErrorContext(c.Ctx, "Failed to create request", "error", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.Logger.ErrorContext(c.Ctx, "Failed to send request", "error", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.Logger.ErrorContext(c.Ctx, "Unexpected status code", "statusCode", resp.StatusCode, "body", string(body))
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var state models.StatesRecord
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		c.Logger.ErrorContext(c.Ctx, "Failed to decode response", "error", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	c.Logger.InfoContext(c.Ctx, "Successfully stopped state record", "runId", runId, "task", string(task), "stateId", state.ID)
	return &state, nil
}
