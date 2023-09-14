package main

import (
	"net/http"
	//nolint:gosec
	_ "net/http/pprof"
	"os"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/options"
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

	options := options.Parse()

	kubeconfig_path := os.Getenv("HOME") + "/.kube/config"
	if options.KubeconfigPath != "" {
		kubeconfig_path = options.KubeconfigPath
	} else if value, present := os.LookupEnv("KUBECONFIG"); present {
		kubeconfig_path = value
	}
	config, _ := clientcmd.
		BuildConfigFromFlags("", kubeconfig_path)
	clientset, _ := kubernetes.NewForConfig(config)

	initial_sync_blacklist, resource_version, err :=
		statistics.CollectInitialPods(options, clientset)
	if err != nil {
		panic(err)
	}

	event_handler := statistics.NewStatisticEventHandler(options, initial_sync_blacklist)

	go event_handler.Run()

	pod_collector := statistics.NewPodCollector(event_handler)
	go pod_collector.Run(clientset, resource_version)

	http.Handle("/metrics", promhttp.Handler())
	handler := zerologhttp.NewHandler(http.DefaultServeMux)
	// No timeouts can be set, but that's OK for us as this HTTP server will not be
	// exposed publicly.
	//nolint:gosec
	if err := http.ListenAndServe(options.ListenAddress, handler); err != nil {
		log.Panic().Err(err).Msg(err.Error())
	}
}
