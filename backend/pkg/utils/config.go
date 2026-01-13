package utils

import "os"

const (
	CONFIG_QUEUE_NAME = "demo"
	CONFIG_IMAGE_NAME = "ghcr.io/lukeelten/micc-creation-day-26-simulator:latest"
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
