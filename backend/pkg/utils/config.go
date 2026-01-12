package utils

import "os"

const (
	CONFIG_QUEUE_NAME = "demo"
)

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

	return "creation-day26"
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
