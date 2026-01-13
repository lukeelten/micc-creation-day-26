package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/lukeelten/micc-creation-day-26/backend/pkg/client"
	"github.com/lukeelten/micc-creation-day-26/backend/pkg/utils"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	runId, exists := os.LookupEnv("RUN_ID")
	if !exists {
		slog.Default().Error("RUN_ID not set, cannot proceed")
		return
	}

	// Config
	backendUrl, exists := os.LookupEnv("BACKEND_URL")
	if !exists {
		slog.Default().Info("BACKEND_URL not set, using default http://localhost:8090/api/")
		backendUrl = "http://localhost:8090/api/"
	}

	var targetDuration time.Duration

	if targetDurationStr, exists := os.LookupEnv("TARGET_DURATION"); exists {
		targetDurationParsed, err := time.ParseDuration(targetDurationStr)
		if err != nil {
			slog.Default().Error("Failed to parse TARGET_DURATION", "error", err)
			return
		}
		targetDuration = targetDurationParsed
	} else {
		targetDuration = utils.RandomDuration()
		slog.Default().Info("TARGET_DURATION not set, using random duration", "duration", targetDuration.String())
	}

	task, ok := os.LookupEnv("RUN_TASK")
	if !ok {
		task = utils.TASK_DOWNLOAD
	} else {
		task = strings.ToLower(task)
	}

	realTask, ok := simuTasks[task]
	if !ok {
		slog.Default().Error("Unknown RUN_TASK", "task", task)
		return
	}

	client := client.NewClient(ctx, backendUrl, slog.Default())
	err := realTask(client, runId, targetDuration)
	if err != nil {
		slog.Default().Error("Task failed", "error", err)
		_, err = client.UpdateRunStatusFailed(runId)
		if err != nil {
			slog.Default().Error("Failed to update run status to failed", "error", err)
		}

		return
	}
}
