package state

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewImagePullStatistic(t *testing.T) {
	configureLogging(t)

	// Define a test pod and container
	container := corev1.Container{
		Name: "test-container",
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "test-namespace",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{container},
		},
	}

	// Call the function
	imagePullStat := NewContainerImagePullStatistic(pod, false, container)

	// Assertions
	assert.Equal(t, pod.Namespace, imagePullStat.podNamespace, "Pod namespace does not match")
	assert.Equal(t, pod.Name, imagePullStat.podName, "Pod name does not match")
	assert.Equal(t, container.Name, imagePullStat.containerName, "Container name does not match")
}

func TestImagePullStatisticLog(t *testing.T) {
	configureLogging(t)

	// Create a buffer to capture the output
	buf := &bytes.Buffer{}

	// Create a test ImagePullStatistic
	now := time.Now()
	imagePullStat := &ContainerImagePullStatistic{
		podNamespace:      "test-namespace",
		podName:           "test-pod",
		containerName:     "test-container",
		alreadyPresent:    true,
		startedTimestamp:  now,
		finishedTimestamp: now.Add(5 * time.Second),
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      imagePullStat.podName,
			Namespace: imagePullStat.podNamespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: imagePullStat.containerName,
				},
			},
		},
	}

	// Call the log function
	imagePullStat.Report(buf, pod, "Test log message")

	// Check if the output contains expected values
	output := buf.String()

	expected := map[string]any{
		"kube_transition_metrics": map[string]any{
			"type": "image_pull",
			"image_pull": map[string]any{
				"already_present":    true,
				"started_timestamp":  imagePullStat.startedTimestamp.Format(time.RFC3339),
				"finished_timestamp": imagePullStat.finishedTimestamp.Format(time.RFC3339),
				// The floating point value here so happens to be the same one we get after marshalling and unmarshalling.
				// Don't touch it, finished timestamp minus started timestamp != 5 seconds by the tinest of margins.
				// https://github.com/stretchr/testify/issues/1576
				"duration_seconds": 5 * time.Second.Seconds(),
			},
			"kube_namespace": "test-namespace",
			"partial":        false,
			"pod_name":       "test-pod",
			"container_name": "test-container",
		},
		"message": "Test log message",
	}

	actual := make(map[string]any)

	err := json.Unmarshal([]byte(output), &actual)
	if assert.NoError(t, err) {
		assert.Equal(t,
			expected["kube_transition_metrics"],
			actual["kube_transition_metrics"],
			"Output does not match expected values")
		assert.Equal(t, expected["message"], actual["message"], "Output does not match expected values")
	}
}
