package statistics

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/options"
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
	statisticEventLoop  *StatisticEventLoop
	options             *options.Options
	imagePullCollectors *sync.Map
}

// NewPodCollector creates a new PodCollector object using the provided
// StatisticEventHandler.
func NewPodCollector(
	options *options.Options,
	eh *StatisticEventLoop,
) *PodCollector {
	return &PodCollector{
		options:             options,
		statisticEventLoop:  eh,
		imagePullCollectors: &sync.Map{},
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

		if _, err := w.statisticEventLoop.PodResync(context.TODO(), resyncUIDs); err != nil {
			log.Panic().Err(err).Msg("Failed to publish resync pods")
		}

		resyncUIDSet := make(map[types.UID]struct{}, len(resyncUIDs))
		for _, uid := range resyncUIDs {
			resyncUIDSet[uid] = struct{}{}
		}

		w.imagePullCollectors.Range(func(key, _ any) bool {
			// We are the only ones using the map, so we can safely cast to types.UID.
			uid, isUID := key.(types.UID)
			if !isUID {
				log.Panic().Any("key", key).Msgf("Non-UID key found in imagePullCollectors map")
			}

			if _, ok := resyncUIDSet[uid]; !ok {
				// Cancel image pull collectors containers who deletion even was missed
				w.cancelImagePullCollector(uid, "pod deleting event missed")
			}

			return true
		})

		w.watch(clientset, resourceVersion)
		log.Warn().Msg("Watch ended, restarting. Some events may be lost.")
		prommetrics.PodCollectorRestarts.Inc()
	}
}

// handlePod processes a Pod event and sends the appropriate statistic event to the statistic event loop.
func (w *PodCollector) handlePod(
	clientset *kubernetes.Clientset,
	eventType watch.EventType,
	pod *corev1.Pod,
) {
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
		w.addImagePullCollector(clientset, pod)

		fallthrough
	case watch.Modified:
		if _, err := w.statisticEventLoop.PodUpdate(context.TODO(), pod); err != nil {
			logger.Error().Err(err).Msg("Error publishing PodUpdate event")
			prommetrics.PodCollectorErrors.Inc()
		}
		if pod.Status.Phase == corev1.PodRunning {
			w.cancelImagePullCollector(pod.UID, "pod already running")
		}
	case watch.Deleted:
		if _, err := w.statisticEventLoop.PodDelete(context.TODO(), pod.UID); err != nil {
			logger.Error().Err(err).Msg("Error publishing PodDelete event")
			prommetrics.PodCollectorErrors.Inc()
		}
		w.cancelImagePullCollector(pod.UID, "pod deleted")
	case watch.Bookmark:
		logger.Warn().Msgf("Got Bookmark event: %+v", pod)
	}
}

// addImagePullCollector adds a new image pull collector for the given pod UID.
// If an image pull collector already exists for the given UID, it is replaced and the old one is cancelled.
func (w *PodCollector) addImagePullCollector(
	clientset *kubernetes.Clientset,
	pod *corev1.Pod,
) {
	collector := newImagePullCollector(w.options, w.statisticEventLoop, pod)
	// Cancel any image pull collectors before removing them from the map
	if existing, ok := w.imagePullCollectors.Swap(pod.UID, collector); ok {
		existingCollector, isCollector := existing.(*imagePullCollector)
		if !isCollector {
			log.Panic().Any("value", existing).Msgf("Non-imagePullCollector found in imagePullCollectors map")
		}
		go existingCollector.cancel("pod replaced")
	}
	go func() {
		collector.Run(clientset)
		// Delete the collector from the map if it is the same as the one that finished running
		w.imagePullCollectors.CompareAndDelete(pod.UID, collector)
	}()
}

// cancelImagePullCollector cancels and removes the image pull collector for the given pod UID.
func (w *PodCollector) cancelImagePullCollector(uid types.UID, reason string) {
	if existing, ok := w.imagePullCollectors.LoadAndDelete(uid); ok {
		collector, ok := existing.(*imagePullCollector)
		if !ok {
			log.Panic().Any("value", existing).Msgf("Non-imagePullCollector found in imagePullCollectors map")
		}
		go collector.cancel(reason)
	}
}

// getWatcher creates a new RetryWatcher for the Pod resource.
// It uses the provided resourceVersion to start watching from that version.
func (w *PodCollector) getWatcher(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	resourceVersion string,
) (*watch_tools.RetryWatcher, error) {
	watcher, err := watch_tools.NewRetryWatcherWithContext(ctx, resourceVersion, &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return clientset.CoreV1().Pods("").List(ctx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return clientset.CoreV1().Pods("").Watch(ctx, options)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return watcher, nil
}

// watch performs the actual watch on the Kubernetes API for all Pod objects.
func (w *PodCollector) watch(
	clientset *kubernetes.Clientset,
	resourceVersion string,
) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	watcher, err := w.getWatcher(ctx, clientset, resourceVersion)
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
		} else {
			w.handlePod(clientset, event.Type, pod)
		}

		prommetrics.PodWatchEvents.With(
			prometheus.Labels{"event_type": string(event.Type)},
		).Inc()
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
	timeOut := w.options.KubeWatchTimeout
	listOptions := metav1.ListOptions{
		TimeoutSeconds: &timeOut,
		Limit:          w.options.KubeWatchMaxEvents,
	}

	blacklistUIDs := make([]types.UID, 0)
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
			clientset.CoreV1().Pods("").List(context.TODO(), listOptions)
		if err != nil {
			log.Error().Err(err).Msg("Error performing initial sync.")

			return nil, "", fmt.Errorf("could not perform initial pod sync: %w", err)
		}

		for _, pod := range list.Items {
			blacklistUIDs = append(blacklistUIDs, pod.UID)
		}
	}
	log.Info().
		Msgf("Initial sync completed, resource version %+v", list.ResourceVersion)

	return blacklistUIDs, list.ResourceVersion, nil
}
