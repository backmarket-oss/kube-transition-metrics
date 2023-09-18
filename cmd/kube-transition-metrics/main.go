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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	prommetrics.Register()

	options := options.Parse()
	zerolog.SetGlobalLevel(options.LogLevel)

	kubeconfigPath := os.Getenv("HOME") + "/.kube/config"
	if options.KubeconfigPath != "" {
		kubeconfigPath = options.KubeconfigPath
	} else if value, present := os.LookupEnv("KUBECONFIG"); present {
		kubeconfigPath = value
	}
	config, _ := clientcmd.
		BuildConfigFromFlags("", kubeconfigPath)
	clientset, _ := kubernetes.NewForConfig(config)

	initialSyncBlacklist, resourceVersion, err :=
		statistics.CollectInitialPods(options, clientset)
	if err != nil {
		panic(err)
	}

	eventHandler := statistics.NewStatisticEventHandler(options, initialSyncBlacklist)

	go eventHandler.Run()

	podCollector := statistics.NewPodCollector(eventHandler)
	go podCollector.Run(clientset, resourceVersion)

	http.Handle("/metrics", promhttp.Handler())
	handler := zerologhttp.NewHandler(http.DefaultServeMux)
	// No timeouts can be set, but that's OK for us as this HTTP server will not be
	// exposed publicly.
	//nolint:gosec
	if err := http.ListenAndServe(options.ListenAddress, handler); err != nil {
		log.Panic().Err(err).Msg(err.Error())
	}
}
