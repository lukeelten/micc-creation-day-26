package utils

import (
	"os"
	"strings"
)

const (
	CONFIG_QUEUE_NAME = "demo"
	CONFIG_IMAGE_NAME = "ghcr.io/lukeelten/micc-creation-day-26:latest"
	DEFAULT_NAMESPACE = "creation-day26"
)

func PtrBool(b bool) *bool {
	return &b
}

func PtrString(s string) *string {
	return &s
}

func PtrInt32(i int32) *int32 {
	return &i
}

func PtrInt64(i int64) *int64 {
	return &i
}

func GetNamespace() string {
	if ns, exists := os.LookupEnv("POD_NAMESPACE"); exists {
		return ns
	}

	if _, exists := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); exists == nil {
		data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
		if err == nil {
			return string(data)
		}
	}

	return DEFAULT_NAMESPACE
}

func GetPodname() string {
	if podname, exists := os.LookupEnv("HOSTNAME"); exists {
		return podname
	}

	if hostname, err := os.Hostname(); err == nil {
		return hostname
	}

	return "unknown-pod"
}

func GetClientBaseUrl() string {
	if internalUrl, exists := os.LookupEnv("INTERNAL_URL"); exists {
		return internalUrl
	}

	return "http://localhost:8080/v1"
}

type AppConfig struct {
	BackendUrl string `json:"backendUrl"`
	Production bool   `json:"production"`
}

func GetAppConfig() AppConfig {
	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	backendUrl := "http://localhost:8090"
	if prodUrl, ok := os.LookupEnv("EXTERNAL_URL"); ok && len(prodUrl) > 0 {
		backendUrl = prodUrl
	}

	config := AppConfig{
		BackendUrl: backendUrl,
		Production: !isGoRun,
	}

	return config
}
