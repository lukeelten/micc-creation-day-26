package utils

import (
	"time"

	"github.com/lukeelten/micc-creation-day-26/backend/pkg/models"
)

const (
	StoreClient  = "k8sClient"
	StoreContext = "globalContext"
)

const (
	TASK_DOWNLOAD = string(models.StatesTaskDownload)
	TASK_CONVERT  = string(models.StatesTaskConvert)
	TASK_PROCESS  = string(models.StatesTaskProcess)
	TASK_UPLOAD   = string(models.StatesTaskUpload)

	TASK_DOWNLOAD_MAX_DURATION = 20 * time.Second
	TASK_CONVERT_MAX_DURATION  = 30 * time.Second
	TASK_PROCESS_MAX_DURATION  = 60 * time.Second
	TASK_UPLOAD_MAX_DURATION   = 15 * time.Second
)
