package state

import (
	"io"
	"iter"
	"time"

	"github.com/Izzette/go-safeconcurrency/eventloop/snapshot"
	"github.com/benbjohnson/immutable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
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

// Containers return an iterator over the containers in the pod.
func (s *PodImagePullStatistic) Containers() iter.Seq2[string, *ContainerImagePullStatistic] {
	return s.EachContainer
}

// EachContainer is an [iter.Seq2] for container names (string) and their corresponding statistics
// [*ContainerImagePullStatistic].
func (s *PodImagePullStatistic) EachContainer(yield func(string, *ContainerImagePullStatistic) bool) {
	containers := s.containers.Iterator()
	for !containers.Done() {
		containerName, container, ok := containers.Next()
		if !ok {
			// This should never happen as we're checking `.Done()` on the iterator.
			log.Panic().Msg("container image pull statistics not found")
		}

		if !yield(containerName, container) {
			break
		}
	}
}

// MapContainers applies the provided function to each container in the pod and returns a new instance of the
// [*PodImagePullStatistic] with the updated container statistics.
// The provided function should return a new instance of the ContainerImagePullStatistic and a boolean indicating
// whether to continue iterating over the containers.
func (s *PodImagePullStatistic) MapContainers(
	apply func(string, *ContainerImagePullStatistic) (*ContainerImagePullStatistic, bool),
) *PodImagePullStatistic {
	// Shallow copy the map to avoid modifying the original instance.
	s = s.Copy()

	for containerName, container := range s.Containers() {
		newContainer, cont := apply(containerName, container)
		if newContainer != container {
			s.containers = s.containers.Set(containerName, newContainer)
		}

		if !cont {
			break
		}
	}

	return s
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

// Partial indicates if the image pull statistic does not contain all the metrics for a complete image pull lifecycle.
func (s *ContainerImagePullStatistic) Partial() bool {
	return s.startedTimestamp.IsZero() || s.finishedTimestamp.IsZero()
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
		Bool("partial", s.Partial()).
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

// ImagePullStatistics holds the statistics for image pulls.
type ImagePullStatistics struct {
	*immutable.Map[types.UID, *PodImagePullStatistic]
}

// NewImagePullStatistics creates a new ImagePullStatistics instance.
func NewImagePullStatistics() *ImagePullStatistics {
	return &ImagePullStatistics{
		Map: immutable.NewMap[types.UID, *PodImagePullStatistic](nil),
	}
}

// Copy implements [github.com/Izzette/go-safeconcurrency/types.Copyable.Copy].
func (s *ImagePullStatistics) Copy() *ImagePullStatistics {
	return snapshot.CopyPtr(s)
}

// Get returns the image pull statistic for the pod with the given UID, if it exists.
func (s *ImagePullStatistics) Get(uid types.UID) (*PodImagePullStatistic, bool) {
	return s.Map.Get(uid)
}

// Set updates the image pull statistic for the pod with the given UID, or adds it if it doesn't exist.
// Set returns a new instance of the ImagePullStatistics with the updated fields.
func (s *ImagePullStatistics) Set(
	uid types.UID,
	imagePullStatistic *PodImagePullStatistic,
) *ImagePullStatistics {
	s = s.Copy()
	s.Map = s.Map.Set(uid, imagePullStatistic)

	return s
}

// Delete deletes the image pull statistic for the pod with the given UID, if it exists.
// Delete returns a new instance of the ImagePullStatistics with the updated fields.
func (s *ImagePullStatistics) Delete(uid types.UID) *ImagePullStatistics {
	s = s.Copy()
	s.Map = s.Map.Delete(uid)

	return s
}

// Len returns the number of image pull statistics.
func (s *ImagePullStatistics) Len() int {
	return s.Map.Len()
}
