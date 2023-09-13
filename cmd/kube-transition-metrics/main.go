package main

import (
	"net/http"
	//nolint:gosec
	_ "net/http/pprof"
	"os"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/prommetrics"
	"github.com/BackMarket-oss/kube-transition-metrics/internal/statistics"
	"github.com/BackMarket-oss/kube-transition-metrics/internal/zerologhttp"
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

	initial_sync_blacklist, resource_version, err :=
		statistics.CollectInitialPods(clientset)
	if err != nil {
		panic(err)
	}

	event_handler := statistics.NewStatisticEventHandler(initial_sync_blacklist)

	go event_handler.Run()

	pod_collector := statistics.NewPodCollector(event_handler)
	go pod_collector.Run(clientset, resource_version)

	http.Handle("/metrics", promhttp.Handler())
	handler := zerologhttp.NewZerologHTTPHandler(http.DefaultServeMux)
	// No timeouts can be set, but that's OK for us as this HTTP server will not be
	// exposed publicly.
	//nolint:gosec
	if err := http.ListenAndServe("0.0.0.0:8080", handler); err != nil {
		log.Panic().Err(err).Msg(err.Error())
	}
}
