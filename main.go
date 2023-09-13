package main

import (
	"net/http"
	//nolint:gosec
	_ "net/http/pprof"
	"os"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/prommetrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	prommetrics.Register()

	config, _ := clientcmd.
		BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	clientset, _ := kubernetes.NewForConfig(config)

	event_handler := NewStatisticEventHandler()
	initial_sync_blacklist, resource_version, err := CollectInitialPods(clientset)
	if err != nil {
		panic(err)
	}

	http.Handle("/metrics", promhttp.Handler())

	event_handler.blacklistUIDs = initial_sync_blacklist
	go event_handler.Run()

	pod_collector := &PodCollector{
		eh: event_handler,
	}
	go pod_collector.Run(clientset, resource_version)

	handler := NewZerologHTTPHandler(http.DefaultServeMux)
	// No timeouts can be set, but that's OK for us.
	//nolint:gosec
	if err := http.ListenAndServe("0.0.0.0:8080", handler); err != nil {
		log.Panic().Err(err).Msg(err.Error())
	}
}
