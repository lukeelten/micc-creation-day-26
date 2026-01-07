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

	app, err := pkg.NewApplication()
	if err != nil {
		panic(err)
	}

	err = app.Run(ctx)

	if err != nil {
		panic(err)
	}
}
