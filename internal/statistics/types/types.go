package types

import (
	"context"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/statistics/state"
	safeconcurrencytypes "github.com/Izzette/go-safeconcurrency/api/types"
	"github.com/rs/zerolog"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// PodCollector is an interface that defines the methods for collecting pod events.
// It is used to allow mocking in tests and to provide a clear contract for the collector's behavior.
//
// Implemented by podCollector in [github.com/BackMarket-oss/kube-transition-metrics/internal/statistics].
type PodCollector interface {
	Run(clientset *kubernetes.Clientset)
}

// ImagePullCollector is an interface that defines the methods for collecting image pull events for a pod.
// It is used to allow mocking in tests and to provide a clear contract for the collector's behavior.
//
// Implemented by imagePullCollector in [github.com/BackMarket-oss/kube-transition-metrics/internal/statistics].
type ImagePullCollector interface {
	Run(clientset *kubernetes.Clientset)
	HandleWatchEvent(watchEvent watch.Event) bool
	HandleEvent(eventType watch.EventType, event *corev1.Event)
	Watch(clientset *kubernetes.Clientset) bool
	WatchOptions() metav1.ListOptions
	Cancel(reason string)
	Logger() *zerolog.Logger
}

// PodStatisticEventLoop is an interface for the pod statistic event loop.
//
// Implemented by podStatisticEventLoop in [github.com/BackMarket-oss/kube-transition-metrics/internal/statistics].
type PodStatisticEventLoop interface {
	safeconcurrencytypes.EventLoop[*state.PodStatistics]

	PodUpdate(ctx context.Context, pod *corev1.Pod) (safeconcurrencytypes.GenerationID, error)
	PodDelete(ctx context.Context, pod *corev1.Pod) (safeconcurrencytypes.GenerationID, error)
	PodResync(ctx context.Context, blacklistUIDs []apimachinerytypes.UID) (safeconcurrencytypes.GenerationID, error)
}

// ImagePullStatisticEventLoop is an interface for the image pull statistic event loop.
//
// Implemented by imagePullStatisticEventLoop in
// [github.com/BackMarket-oss/kube-transition-metrics/internal/statistics].
type ImagePullStatisticEventLoop interface {
	safeconcurrencytypes.EventLoop[*state.ImagePullStatistics]
	ImagePullUpdate(
		ctx context.Context, pod *corev1.Pod, k8sEvent *corev1.Event) (safeconcurrencytypes.GenerationID, error)
	ImagePullDelete(ctx context.Context, pod *corev1.Pod) (safeconcurrencytypes.GenerationID, error)
}
