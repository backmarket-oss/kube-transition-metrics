package state

import (
	"io"
	"net/url"
	"path"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/util/parsers"
)

// appLabelFields is a map of metric labels to Kubernetes application labels.
// It is used by [appLabels] to add application labels to the event.
//
//nolint:gochecknoglobals // This is a constant map of metric labels to Kubernetes application labels.
var appLabelFields = map[string]string{
	"kube_app_component":  "app.kubernetes.io/component",
	"kube_app_instance":   "app.kubernetes.io/instance",
	"kube_app_managed_by": "app.kubernetes.io/managed-by",
	"kube_app_name":       "app.kubernetes.io/name",
	"kube_app_part_of":    "app.kubernetes.io/part-of",
	"kube_app_version":    "app.kubernetes.io/version",
}

// commonPodLabels returns a function that adds common pod labels to the event.
// This can be used with [zerolog.Event.Func] to add labels to the event.
func commonPodLabels(pod *corev1.Pod) func(event *zerolog.Event) {
	return func(event *zerolog.Event) {
		event.Str("kube_namespace", pod.Namespace)
		event.Str("pod_name", pod.Name)
		if pod.Spec.NodeName != "" {
			event.Str("kube_node", pod.Spec.NodeName)
		}
		if pod.Status.QOSClass != "" {
			event.Str("kube_qos", string(pod.Status.QOSClass))
		}
		if pod.Spec.PriorityClassName != "" {
			event.Str("kube_priority_class", pod.Spec.PriorityClassName)
		}
		if pod.Spec.RuntimeClassName != nil {
			event.Str("kube_runtime_class", *pod.Spec.RuntimeClassName)
		}
		event.Func(ownerRefLabels(pod.OwnerReferences))
		event.Func(appLabels(pod.Labels))
		// TODO: find Service pointing to this pod and add kube_service label.
		// TODO: add custom labels from commandline flag.
	}
}

// commonContainerLabels returns a function that adds common container labels to the event.
// This can be used with [zerolog.Event.Func] to add labels to the event.
func commonContainerLabels(logger *zerolog.Logger, container *corev1.Container) func(event *zerolog.Event) {
	return func(event *zerolog.Event) {
		event.Str("container_name", container.Name)
		event.Func(imageLabels(logger, container.Image))
	}
}

// imageLabels returns a function that adds image labels to the event.
// This can be used with [zerolog.Event.Func] to add labels to the event.
func imageLabels(logger *zerolog.Logger, image string) func(event *zerolog.Event) {
	return func(event *zerolog.Event) {
		repo, tag, digest, err := parsers.ParseImageName(image)
		if err != nil {
			logger.Error().Err(err).Str("image", image).Msg("failed to parse image name")

			return
		}

		event.Str("image_name", repo)
		parsed, err := url.Parse(repo)
		if err != nil {
			logger.Error().Err(err).Str("image_repo", repo).Msg("failed to parse image repo")
		} else {
			shortImage := path.Base(parsed.Path)
			event.Str("short_image", shortImage)
		}

		if tag == "" {
			tag = digest
		}
		event.Str("image_tag", tag)
	}
}

// ownerRefLabels returns a function that adds owner reference labels to the event.
// This can be used with [zerolog.Event.Func] to add labels to the event.
func ownerRefLabels(ownerRefs []metav1.OwnerReference) func(event *zerolog.Event) {
	return func(event *zerolog.Event) {
		ownerRef := controllerRef(ownerRefs)
		if ownerRef == nil {
			return
		}

		event.Str("kube_ownerref_kind", strings.ToLower(ownerRef.Kind))
		event.Str("kube_ownerref_name", ownerRef.Name)

		// TODO: recursively resolve cronjob and deployment as well.
		switch ownerRef.Kind {
		case "DaemonSet":
			event.Str("kube_daemon_set", ownerRef.Name)
		case "Job":
			event.Str("kube_job", ownerRef.Name)
		case "ReplicaSet":
			event.Str("kube_replica_set", ownerRef.Name)
		case "StatefulSet":
			event.Str("kube_stateful_set", ownerRef.Name)
		}
	}
}

// appLabels returns a function that adds application labels to the event.
// This can be used with [zerolog.Event.Func] to add labels to the event.
func appLabels(labels map[string]string) func(event *zerolog.Event) {
	return func(event *zerolog.Event) {
		for metricLabel, k8sLabel := range appLabelFields {
			if value, ok := labels[k8sLabel]; ok {
				event.Str(metricLabel, value)
			}
		}
	}
}

// controllerRef returns the controller owner reference from the list of owner references.
// If no controller is found, it returns nil.
func controllerRef(ownerRefs []metav1.OwnerReference) *metav1.OwnerReference {
	for _, ownerRef := range ownerRefs {
		if ownerRef.Controller != nil && *ownerRef.Controller {
			return &ownerRef
		}
	}

	return nil
}

// logMetrics writes the kube transition metrics to the output writer.
func logMetrics(output io.Writer, metricType string, metrics *zerolog.Event, message string) {
	logger := log.Output(output)
	logger.Log().
		Dict("kube_transition_metrics", metrics.Str("type", metricType)).
		Msg(message)
}

// findContainer finds the container from the list by specified name.
// It returns nil if the container is not found.
func findContainer(name string, containers []corev1.Container) *corev1.Container {
	for _, container := range containers {
		if container.Name == name {
			return &container
		}
	}

	return nil
}
