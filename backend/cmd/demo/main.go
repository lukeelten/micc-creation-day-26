package main

import (
	"context"
	"flag"
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

	var runId string
	var backendUrl string
	var targetDurationStr string
	var task string

	flag.StringVar(&runId, "run-id", os.Getenv("RUN_ID"), "Run ID (required)")
	flag.StringVar(&backendUrl, "backend-url", utils.GetClientBaseUrl(), "Backend URL")
	flag.StringVar(&targetDurationStr, "target-duration", "", "Target duration (e.g. 5s, 1m)")
	flag.StringVar(&task, "task", utils.TASK_DOWNLOAD, "Task to run")
	flag.Parse()

	if runId == "" {
		slog.Default().Error("run-id is required")
		os.Exit(1)
	}

	var targetDuration time.Duration

	if targetDurationStr != "" {
		targetDurationParsed, err := time.ParseDuration(targetDurationStr)
		if err != nil {
			slog.Default().Error("Failed to parse target-duration", "error", err)
			os.Exit(2)
		}

		targetDuration = targetDurationParsed
	} else {
		targetDuration = utils.RandomDuration()
		slog.Default().Info("target-duration not set, using random duration", "duration", targetDuration.String())
	}

	task = strings.ToLower(task)

	realTask, ok := simuTasks[task]
	if !ok {
		slog.Default().Error("Unknown task", "task", task)
		os.Exit(3)
		return
	}

	slog.Default().Info("Starting task execution",
		"run-id", runId,
		"backend-url", backendUrl,
		"target-duration", targetDuration.String(),
		"task", task,
	)

	client := client.NewClient(ctx, backendUrl, slog.Default())
	err := client.StartState(runId, task)
	if err != nil {
		slog.Default().Error("Failed to start state", "error", err)
		os.Exit(4)
	} else {
		slog.Default().Info("State started successfully", "task", task)
	}

	err = realTask(client, runId, targetDuration)
	if err != nil {
		slog.Default().Error("Task failed", "error", err)
		_, err = client.UpdateRunStatusFailed(runId)
		if err != nil {
			slog.Default().Error("Failed to update run status to failed", "error", err)
		}

		os.Exit(5)
	}

	err = client.StopState(runId, task)
	if err != nil {
		slog.Default().Error("Failed to stop state", "error", err)
		os.Exit(6)
	}

	slog.Default().Info("Task completed successfully", "task", task)
}
