package controller

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/lukeelten/micc-creation-day-26/backend/pkg/models"
	"github.com/lukeelten/micc-creation-day-26/backend/pkg/utils"
	"github.com/pocketbase/pocketbase"
	"k8s.io/client-go/kubernetes"
)

type RunController struct {
	Logger *slog.Logger
	Client *kubernetes.Clientset
	Pb     *pocketbase.PocketBase
	Ctx    context.Context
	wg     sync.WaitGroup
}

func NewRunController(pb *pocketbase.PocketBase) (*RunController, error) {
	k8sClient, ok := pb.Store().Get(utils.StoreClient).(*kubernetes.Clientset)
	if !ok {
		return nil, errors.New("Invalid kubernetes client found")
	}

	return &RunController{
		Logger: pb.Logger(),
		Client: k8sClient,
		Pb:     pb,
		Ctx:    pb.RootCmd.Context(),
	}, nil
}

func (rc *RunController) Start() error {
	go func() {
		// Wait two seconds to allow the application to fully start
		time.Sleep(2 * time.Second)
		rc.Logger.Info("RunController started")

		rc.Bootstrap()

		rc.wg.Wait()
		rc.Logger.Info("RunController stopped")
	}()

	return nil
}

func (rc *RunController) Bootstrap() error {
	// Fetch all runs
	allRuns, err := rc.Pb.FindAllRecords(models.CollectionsRuns)
	if err != nil {
		return err
	}

	for _, r := range allRuns {
		runRecord := models.ConvertRunRecord(r)

		runInstance := NewRunInstance(rc, runRecord)
		rc.wg.Go(runInstance.Start)
	}

	return nil
}
