package pkg

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// CreateKubernetesJob creates a simple Job for the provided taskName.
// It uses env vars JOB_IMAGE (default: busybox) and JOB_NAMESPACE (default: default).
func CreateKubernetesJob(ctx context.Context, logger *slog.Logger, clientset *kubernetes.Clientset, taskName string) (*batchv1.Job, error) {
	image := os.Getenv("JOB_IMAGE")
	if image == "" {
		image = "busybox"
	}

	namespace := os.Getenv("JOB_NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	logger.InfoContext(ctx, "Creating Kubernetes Job", "task", taskName, "image", image, "namespace", namespace)

	backoffLimit := int32(0)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "task-" + taskName + "-",
			Labels: map[string]string{
				"app":        "micc-creation-day-26",
				"task":       taskName,
			},
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &backoffLimit,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  "micc-creation-day-26",
						"task": taskName,
					},
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  "task",
							Image: image,
							Command: []string{"sh", "-c"},
							Args: []string{"echo Running task: $TASK_NAME; sleep 1"},
							Env: []corev1.EnvVar{
								{Name: "TASK_NAME", Value: taskName},
							},
						},
					},
				},
			},
		},
	}

	created, err := clientset.BatchV1().Jobs(namespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		logger.ErrorContext(ctx, "Failed to create Job", "error", err)
		return nil, err
	}

	logger.InfoContext(ctx, "Job created", "name", created.Name, "namespace", namespace)
	return created, nil
}
