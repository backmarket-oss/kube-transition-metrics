package state

import (
	"io"
	"time"

	"github.com/Izzette/go-safeconcurrency/eventloop/snapshot"
	"github.com/benbjohnson/immutable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
)

// PodImagePullStatistic holds the statistics for a pod image pull.
// PodImagePullStatistic is immutable, all the methods return a new instance of the struct.
// Do not lose track of the returned instance, it should be assigned to the containing structure.
type PodImagePullStatistic struct {
	podNamespace string
	podName      string

	containers *immutable.Map[string, *ContainerImagePullStatistic]
}

// NewPodImagePullStatistic creates a new PodImagePullStatistic instance populated with the containers in the pod.
func NewPodImagePullStatistic(pod *corev1.Pod) *PodImagePullStatistic {
	containers := immutable.NewMapBuilder[string, *ContainerImagePullStatistic](nil)
	for _, c := range pod.Spec.InitContainers {
		containers.Set(c.Name, NewContainerImagePullStatistic(pod, true, c))
	}
	for _, c := range pod.Spec.Containers {
		containers.Set(c.Name, NewContainerImagePullStatistic(pod, false, c))
	}

	return &PodImagePullStatistic{
		podNamespace: pod.Namespace,
		podName:      pod.Name,
		containers:   containers.Map(),
	}
}

// Copy implements [github.com/Izzette/go-safeconcurrency/types.Copyable.Copy].
func (s *PodImagePullStatistic) Copy() *PodImagePullStatistic {
	return snapshot.CopyPtr(s)
}

// Get returns the image pull statistic for the container with the given name, if it exists.
func (s *PodImagePullStatistic) Get(containerName string) (*ContainerImagePullStatistic, bool) {
	return s.containers.Get(containerName)
}

// Set updates the image pull statistic for the container with the given name, or adds it if it doesn't exist.
// Set returns a new instance of the PodImagePullStatistic with the updated fields.
func (s *PodImagePullStatistic) Set(
	container *ContainerImagePullStatistic,
) *PodImagePullStatistic {
	s = s.Copy()
	s.containers = s.containers.Set(container.containerName, container)

	return s
}

// ContainerImagePullStatistic holds the statistics for a container image pull.
// ContainerImagePullStatistic is immutable, all the methods return a new instance of the struct.
// Do not lose track of the returned instance, it should be assigned to the containing structure.
type ContainerImagePullStatistic struct {
	podNamespace  string
	podName       string
	containerName string
	initContainer bool

	alreadyPresent    bool
	startedTimestamp  time.Time
	finishedTimestamp time.Time
}

// NewContainerImagePullStatistic creates a new ContainerImagePullStatistic instance.
// The provided pod and container are used to populate the podNamespace, podName, and containerName fields.
func NewContainerImagePullStatistic(
	pod *corev1.Pod,
	initContainer bool,
	container corev1.Container,
) *ContainerImagePullStatistic {
	return &ContainerImagePullStatistic{
		podNamespace:  pod.Namespace,
		podName:       pod.Name,
		containerName: container.Name,
		initContainer: initContainer,
	}
}

// Copy implements [github.com/Izzette/go-safeconcurrency/types.Copyable.Copy].
func (s *ContainerImagePullStatistic) Copy() *ContainerImagePullStatistic {
	return snapshot.CopyPtr(s)
}

// Update updates the image pull statistic with the provided event.
// If the event is a pull event, it sets the startedTimestamp.
func (s *ContainerImagePullStatistic) Update(event *corev1.Event) *ContainerImagePullStatistic {
	s = s.Copy()

	switch event.Reason {
	case "Pulled":
		if s.finishedTimestamp.IsZero() {
			s.finishedTimestamp = event.LastTimestamp.Time
		}
		// If we never received a Pulling event, we assume the image was already present.
		if s.startedTimestamp.IsZero() {
			s.alreadyPresent = true
		}

		fallthrough
	case "Pulling":
		if s.startedTimestamp.IsZero() {
			s.startedTimestamp = event.LastTimestamp.Time
		}
	}

	return s
}

// Report logs the image pull statistic to the provided output writer.
func (s *ContainerImagePullStatistic) Report(output io.Writer, pod *corev1.Pod, message string) {
	logger := s.logger()

	var container *corev1.Container
	if s.initContainer {
		container = findContainer(s.containerName, pod.Spec.InitContainers)
	} else {
		container = findContainer(s.containerName, pod.Spec.Containers)
	}
	if container == nil {
		logger.Panic().Msg("container not found")
	}

	metrics := zerolog.Dict().
		Func(commonPodLabels(pod)).
		Func(commonContainerLabels(&logger, container)).
		Dict("image_pull", s.event())
	logMetrics(output, "image_pull", metrics, message)
}

// event returns a zerolog event with the image pull statistics.
func (s *ContainerImagePullStatistic) event() *zerolog.Event {
	event := zerolog.Dict()
	event.Bool("already_present", s.alreadyPresent)
	if !s.startedTimestamp.IsZero() {
		event.Time("started_timestamp", s.startedTimestamp)
	}
	if !s.finishedTimestamp.IsZero() {
		event.Time("finished_timestamp", s.finishedTimestamp)
		if !s.startedTimestamp.IsZero() {
			event.Dur("duration_seconds", s.finishedTimestamp.Sub(s.startedTimestamp))
		}
	}

	return event
}

// logger returns a logger for the container image pull statistic.
func (s *ContainerImagePullStatistic) logger() zerolog.Logger {
	return log.With().
		Str("kube_namespace", s.podNamespace).
		Str("pod_name", s.podName).
		Str("container_name", s.containerName).
		Logger()
}
