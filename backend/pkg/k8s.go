package pkg

import (
	"log/slog"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func NewKubernetesClient(logger *slog.Logger, kubeconfigFlag *string) (*kubernetes.Clientset, error) {
	if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token"); err == nil {
		logger.Info("Creating InCluster kubernetes client")
		return inClusterClient(logger)
	}

	return externalClient(logger, *kubeconfigFlag)
}

func inClusterClient(logger *slog.Logger) (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Error("Failed to create in-cluster config", "error", err)
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func externalClient(logger *slog.Logger, manualPath string) (*kubernetes.Clientset, error) {
	var path string

	kubeconfigEnv, hasKubeconfigEnv := os.LookupEnv("KUBECONFIG")
	if hasKubeconfigEnv {
		path = kubeconfigEnv
	} else if len(manualPath) > 0 {
		path = manualPath
	} else if _, err := os.Stat("./kubeconfig"); err == nil {
		path = "./kubeconfig"
	} else {
		path = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	logger.Info("Creating external kubernetes client", slog.String("kubeconfigPath", path))

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		logger.Error("Failed to build config from flags", "error", err)
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
