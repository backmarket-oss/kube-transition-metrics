package statistics

import (
	"context"
	"fmt"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/prommetrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type PodCollector struct {
	eh *StatisticEventHandler
}

func NewPodCollector(eh *StatisticEventHandler) *PodCollector {
	return &PodCollector{
		eh: eh,
	}
}

func CollectInitialPods(
	clientset *kubernetes.Clientset,
) ([]types.UID, string, error) {
	time_out := int64(60)
	list_options := metav1.ListOptions{TimeoutSeconds: &time_out, Limit: 100}

	blacklist_uids := make([]types.UID, 0)
	log.Info().Msg("Listing pods to get initial state ...")
	var list *corev1.PodList
	for list == nil || list.Continue != "" {
		if list != nil {
			log.Debug().Msgf("Initial list contains %d items ...", len(list.Items))
			list_options.Continue = list.Continue
		}

		log.Debug().Msgf("Listing from %+v ...", list_options.Continue)
		var err error
		list, err =
			clientset.CoreV1().Pods("").List(context.Background(), list_options)
		if err != nil {
			log.Error().Err(err).Msg("Error performing initial sync.")

			return nil, "", fmt.Errorf("could not perform initial pod sync: %w", err)
		}

		for _, pod := range list.Items {
			blacklist_uids = append(blacklist_uids, pod.UID)
		}
	}
	log.Info().
		Msgf("Initial sync completed, resource version %+v", list.ResourceVersion)

	return blacklist_uids, list.ResourceVersion, nil
}

type PodAddedEvent struct {
	collector *PodCollector
	pod       *corev1.Pod
	clientset *kubernetes.Clientset
}

func (ev *PodAddedEvent) PodUID() types.UID {
	return ev.pod.UID
}

func (ev *PodAddedEvent) Handle(statistic *PodStatistic) bool {
	// As the PodAddedEvent may be called more than once, the initialization must only happen once.
	if !statistic.Initialized {
		statistic.Initialize(ev.pod)
		statistic.ImagePullCollector = NewImagePullCollector(ev.collector.eh, ev.pod.Namespace, ev.pod.UID)
		go statistic.ImagePullCollector.Run(ev.clientset)
	}

	statistic.update(ev.pod)

	return false
}

type PodModifiedEvent struct {
	pod *corev1.Pod
}

func (ev *PodModifiedEvent) PodUID() types.UID {
	return ev.pod.UID
}

func (ev *PodModifiedEvent) Handle(statistic *PodStatistic) bool {
	statistic.update(ev.pod)

	return false
}

type PodDeletedEvent struct {
	uid types.UID
}

func (ev *PodDeletedEvent) PodUID() types.UID {
	return ev.uid
}

func (ev *PodDeletedEvent) Handle(statistic *PodStatistic) bool {
	go statistic.ImagePullCollector.cancel("pod_deleted")

	return true
}

func (w *PodCollector) handlePod(
	clientset *kubernetes.Clientset,
	event_type watch.EventType,
	pod *corev1.Pod,
) StatisticEvent {
	logger := log.With().
		Str("kube_namespace", pod.Namespace).
		Str("pod_name", pod.Name).
		Str("pod_uid", string(pod.UID)).
		Str("event_type", string(event_type)).
		Logger()
	logger.Debug().Msg("Collecting statistics for pod")

	// The watch.EventType watch.Error is already tested in the caller, as if there
	// is an error no pod is sent.
	//nolint:exhaustive
	switch event_type {
	case watch.Added:
		return &PodAddedEvent{
			clientset: clientset,
			collector: w,
			pod:       pod,
		}
	case watch.Modified:
		return &PodModifiedEvent{
			pod: pod,
		}
	case watch.Deleted:
		return &PodDeletedEvent{
			uid: pod.UID,
		}
	case watch.Bookmark:
		logger.Warn().Msgf("Got Bookmark event: %+v", pod)
	}

	return nil
}

func (w *PodCollector) watch(
	clientset *kubernetes.Clientset,
	resource_version string,
) {
	time_out := int64(60)
	send_initial_events := resource_version != ""
	watch_ops := metav1.ListOptions{
		TimeoutSeconds:    &time_out,
		SendInitialEvents: &send_initial_events,
		Watch:             true,
		ResourceVersion:   resource_version,
		Limit:             100,
	}
	watcher, err :=
		clientset.CoreV1().Pods("").Watch(context.Background(), watch_ops)
	if err != nil {
		log.Panic().Err(err).Msg("Error starting watcher.")
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		var pod *corev1.Pod
		var is_a_pod bool
		if event.Type == watch.Error {
			log.Error().Msgf("Watch event error: %+v", event)
			prommetrics.POD_COLLECTOR_ERRORS.Inc()

			break
		} else if pod, is_a_pod = event.Object.(*corev1.Pod); !is_a_pod {
			log.Panic().Msgf("Watch event is not a Pod: %+v", event)
		} else if event := w.handlePod(clientset, event.Type, pod); event != nil {
			w.eh.EventChan <- event
		}

		prommetrics.PODS_PROCESSED.With(prometheus.Labels{"event_type": string(event.Type)}).Inc()
	}
}

func (w *PodCollector) Run(
	clientset *kubernetes.Clientset,
	resource_version string,
) {
	for {
		w.watch(clientset, resource_version)

		// Some leak in w.blacklistUids and w.statistics could happen, as Deleted
		// events may be lost. This could be mitigated by performing another full List
		// and checking for removed pod UIDs.
		log.Warn().Msg("Watch ended, restarting. Some events may be lost.")
		prommetrics.POD_COLLECTOR_RESTARTS.Inc()
		resource_version = ""
	}
}
