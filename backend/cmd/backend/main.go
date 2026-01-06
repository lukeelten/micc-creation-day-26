package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lukeelten/micc-creation-day-26/backend/pkg"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	app := pkg.NewApplication()

	err := app.Run(ctx)

	if err != nil {
		panic(err)
	}
}
