package controller

import (
	"fmt"
	"os"
	"time"

	"github.com/lukeelten/micc-creation-day-26/backend/pkg/utils"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateKubernetesJob creates a simple Job for the provided taskName.
func (rc *RunController) CreateKubernetesJob(runId string, taskName string, taskDuration time.Duration) (*batchv1.Job, error) {
	rc.Logger.Debug("creating kubernetes job with parameters", "runId", runId, "taskName", taskName, "taskDuration", taskDuration.String())

	namespace := os.Getenv("JOB_NAMESPACE")
	if namespace == "" {
		namespace = utils.GetNamespace()
	}

	jobArgs := []string{
		"-backend-url", utils.GetClientBaseUrl(),
		"-run-id", runId,
		"-task", taskName,
		"-target-duration", taskDuration.String(),
	}

	maxRuntime := taskDuration + (5 * time.Minute)
	ttlAfterFinished := 1 * time.Hour

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("run-%s-%s-", runId, taskName),
			Labels: map[string]string{
				"app":                       "micc-creation-day-26",
				"task":                      taskName,
				"runId":                     runId,
				"kueue.x-k8s.io/queue-name": utils.CONFIG_QUEUE_NAME,
			},
		},
		Spec: batchv1.JobSpec{
			BackoffLimit:            utils.PtrInt32(5),
			Parallelism:             utils.PtrInt32(1),
			ActiveDeadlineSeconds:   utils.PtrInt64(int64(maxRuntime.Seconds())),
			TTLSecondsAfterFinished: utils.PtrInt32(int32(ttlAfterFinished.Seconds())),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":   "micc-creation-day-26",
						"task":  taskName,
						"runId": runId,
					},
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: utils.PtrBool(true),
					},
					Containers: []corev1.Container{
						{
							Name:    "task",
							Image:   utils.CONFIG_IMAGE_NAME,
							Command: []string{"/app/demo"},
							Args:    jobArgs,
							SecurityContext: &corev1.SecurityContext{
								Privileged:               utils.PtrBool(false),
								ReadOnlyRootFilesystem:   utils.PtrBool(true),
								AllowPrivilegeEscalation: utils.PtrBool(false),
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("50m"),
									corev1.ResourceMemory: resource.MustParse("64Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("256Mi"),
								},
							},
						},
					},
				},
			},
		},
	}

	created, err := rc.Client.BatchV1().Jobs(namespace).Create(rc.Ctx, job, metav1.CreateOptions{})
	if err != nil {
		rc.Logger.ErrorContext(rc.Ctx, "Failed to create Job", "error", err)
		return nil, err
	}

	rc.Logger.InfoContext(rc.Ctx, "Job created", "name", created.Name, "namespace", namespace)
	return created, nil
}
