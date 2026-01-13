package utils

import "github.com/lukeelten/micc-creation-day-26/backend/pkg/models"

const (
	StoreClient  = "k8sClient"
	StoreContext = "globalContext"
)

const (
	TASK_DOWNLOAD = string(models.StatesTaskDownload)
	TASK_CONVERT  = string(models.StatesTaskConvert)
	TASK_PROCESS  = string(models.StatesTaskProcess)
	TASK_UPLOAD   = string(models.StatesTaskUpload)
)
