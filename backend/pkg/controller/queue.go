package controller

func (rc *RunController) StartQueue() {

	rc.wg.Go(func() {
		rc.Logger.Debug("Work queue started")
		defer rc.Logger.Debug("Work queue stopped")

		for {
			select {
			case runId := <-rc.workQueue:
				rc.Logger.Info("Processing run from work queue", "runId", runId)

				err := rc.ProcessRun(runId)
				if err != nil {
					rc.Logger.Error("Error processing run", "runId", runId, "error", err)
				}

				rc.Logger.Debug("Processing finished", "runId", runId)
			case <-rc.Ctx.Done():
				rc.Logger.Info("Work queue stopped")
				return
			}
		}
	})
}
