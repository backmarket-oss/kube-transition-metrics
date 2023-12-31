package statistics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/prommetrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	watch_tools "k8s.io/client-go/tools/watch"
)

// PodCollector uses the Kubernetes Watch API to monitor for all changes on Pods
// and send statistic events to the StatisticEventHandler to track created,
// modified, and deleted Pods during their lifecycles.
type PodCollector struct {
	eh *StatisticEventHandler
}

// NewPodCollector creates a new PodCollector object using the provided
// StatisticEventHandler.
func NewPodCollector(
	eh *StatisticEventHandler,
) *PodCollector {
	return &PodCollector{
		eh: eh,
	}
}

// collectInitialPods generates a list of Pod UIDs currently existing on the
// cluster. This is used to filter pre-existing Pods by the
// StatisticEventHandler to avoid generating inaccurate or incomplete metrics.
// It returns the list of Pod UIDs, the resource version for these UIDs, and an
// error if one occurred.
func (w *PodCollector) collectInitialPods(
	clientset *kubernetes.Clientset,
) ([]types.UID, string, error) {
	timeOut := w.eh.options.KubeWatchTimeout
	listOptions := metav1.ListOptions{
		TimeoutSeconds: &timeOut,
		Limit:          w.eh.options.KubeWatchMaxEvents,
	}

	blacklistUids := make([]types.UID, 0)
	log.Info().Msg("Listing pods to get initial state ...")
	var list *corev1.PodList
	for list == nil || list.Continue != "" {
		if list != nil {
			log.Debug().Msgf("Initial list contains %d items ...", len(list.Items))
			listOptions.Continue = list.Continue
		}

		log.Debug().Msgf("Listing from %+v ...", listOptions.Continue)
		var err error
		list, err =
			clientset.CoreV1().Pods("").List(context.Background(), listOptions)
		if err != nil {
			log.Error().Err(err).Msg("Error performing initial sync.")

			return nil, "", fmt.Errorf("could not perform initial pod sync: %w", err)
		}

		for _, pod := range list.Items {
			blacklistUids = append(blacklistUids, pod.UID)
		}
	}
	log.Info().
		Msgf("Initial sync completed, resource version %+v", list.ResourceVersion)

	return blacklistUids, list.ResourceVersion, nil
}

type podAddedEvent struct {
	collector *PodCollector
	pod       *corev1.Pod
	clientset *kubernetes.Clientset
}

func (ev *podAddedEvent) podUID() types.UID {
	return ev.pod.UID
}

func (ev *podAddedEvent) handle(statistic *podStatistic) bool {
	// As the PodAddedEvent may be called more than once, the initialization must
	// only happen once.
	if !statistic.initialized {
		statistic.initialize(ev.pod)
		statistic.imagePullCollector = newImagePullCollector(
			ev.collector.eh,
			ev.pod.Namespace,
			ev.pod.UID,
		)
		go statistic.imagePullCollector.Run(ev.clientset)
	}

	statistic.update(ev.pod)

	return false
}

type podModifiedEvent struct {
	pod *corev1.Pod
}

func (ev *podModifiedEvent) podUID() types.UID {
	return ev.pod.UID
}

func (ev *podModifiedEvent) handle(statistic *podStatistic) bool {
	statistic.update(ev.pod)

	return false
}

type podDeletedEvent struct {
	uid types.UID
}

func (ev *podDeletedEvent) podUID() types.UID {
	return ev.uid
}

func (ev *podDeletedEvent) handle(statistic *podStatistic) bool {
	go statistic.imagePullCollector.cancel("pod_deleted")

	return true
}

func (w *PodCollector) handlePod(
	clientset *kubernetes.Clientset,
	eventType watch.EventType,
	pod *corev1.Pod,
) statisticEvent {
	logger := log.With().
		Str("kube_namespace", pod.Namespace).
		Str("pod_name", pod.Name).
		Str("pod_uid", string(pod.UID)).
		Str("event_type", string(eventType)).
		Logger()
	logger.Debug().Msg("Collecting statistics for pod")

	// The watch.EventType watch.Error is already tested in the caller, as if there
	// is an error no pod is sent.
	//nolint:exhaustive
	switch eventType {
	case watch.Added:
		return &podAddedEvent{
			clientset: clientset,
			collector: w,
			pod:       pod,
		}
	case watch.Modified:
		return &podModifiedEvent{
			pod: pod,
		}
	case watch.Deleted:
		return &podDeletedEvent{
			uid: pod.UID,
		}
	case watch.Bookmark:
		logger.Warn().Msgf("Got Bookmark event: %+v", pod)
	}

	return nil
}

func (w *PodCollector) getWatcher(
	clientset *kubernetes.Clientset,
	resourceVersion string,
) (*watch_tools.RetryWatcher, error) {
	watcher, err := watch_tools.NewRetryWatcher(resourceVersion, &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			//nolint:wrapcheck
			return clientset.CoreV1().Pods("").List(context.Background(), options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			//nolint:wrapcheck
			return clientset.CoreV1().Pods("").Watch(context.Background(), options)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return watcher, nil
}

func (w *PodCollector) watch(
	clientset *kubernetes.Clientset,
	resourceVersion string,
) {
	watcher, err := w.getWatcher(clientset, resourceVersion)
	if err != nil {
		log.Panic().Err(err).Msg("Error starting watcher.")
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		var pod *corev1.Pod
		var isAPod bool
		if event.Type == watch.Error {
			prommetrics.PodCollectorErrors.Inc()
			// API error will not be wrapped and StatusError doesn't implement the
			// nessesary interface.
			//nolint:errorlint
			apiStatus, ok := apierrors.FromObject(event.Object).(*apierrors.StatusError)
			if ok && apiStatus.ErrStatus.Code == http.StatusGone {
				// The resource version we were watching is too old.
				log.Warn().Msg("Resource version too old, resetting watch.")

				return
			} else {
				log.Error().Msgf("Watch event error: %+v", event)
			}

			continue // RetryWatcher will handle reconnection, so just continue
		} else if pod, isAPod = event.Object.(*corev1.Pod); !isAPod {
			log.Panic().Msgf("Watch event is not a Pod: %+v", event)
		} else if event := w.handlePod(clientset, event.Type, pod); event != nil {
			w.eh.Publish(event)
		}

		prommetrics.PodsProcessed.With(
			prometheus.Labels{"event_type": string(event.Type)},
		).Inc()
	}
}

// Run watches the Kubernetes Pods objects and reports them to the
// StatisticEventHandler used to initialize the PodCollector. It is blocking and
// should be run in another goroutine to the StatisticEventHandler and other
// collectors.
func (w *PodCollector) Run(clientset *kubernetes.Clientset) {
	for {
		resyncUIDs, resourceVersion, err := w.collectInitialPods(clientset)
		if err != nil {
			log.Panic().Err(err).Msg(
				"Failed to resync after 410 Gone from kubernetes Watch API")
		}

		w.eh.resyncChan.Publish(resyncUIDs)
		w.watch(clientset, resourceVersion)
		log.Warn().Msg("Watch ended, restarting. Some events may be lost.")
		prommetrics.PodCollectorRestarts.Inc()
	}
}
