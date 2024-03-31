package main

import (
	"net/http"
	//nolint:gosec
	_ "net/http/pprof"
	"os"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/logging"
	"github.com/BackMarket-oss/kube-transition-metrics/internal/options"
	"github.com/BackMarket-oss/kube-transition-metrics/internal/prommetrics"
	"github.com/BackMarket-oss/kube-transition-metrics/internal/statistics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getKubeconfigInCluster() *rest.Config {
	if _, present := os.LookupEnv("KUBERNETES_SERVICE_HOST"); !present {
		return nil
	}

	// We're inside a pod, use in-cluster configuration
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to get in-cluster Kubernetes config")
	}

	return config
}

func getKubeconfigFromPath(options *options.Options) *rest.Config {
	var kubeconfigPath string
	if options.KubeconfigPath != "" {
		kubeconfigPath = options.KubeconfigPath
	} else if value, present := os.LookupEnv("KUBECONFIG"); present {
		kubeconfigPath = value
	} else {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Panic().Err(err).
			Msgf("Failed to build Kubernetes config from %s", kubeconfigPath)
	}

	return config
}

func getKubeconfig(options *options.Options) *rest.Config {
	if options.KubeconfigPath == "" {
		if config := getKubeconfigInCluster(); config != nil {
			return config
		}
	}

	return getKubeconfigFromPath(options)
}

func main() {
	logging.Configure()
	prommetrics.Register()

	options := options.Parse()
	logging.SetOptions(options)

	config := getKubeconfig(options)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to build kubernetes client")
	}

	eventHandler := statistics.NewStatisticEventHandler(options)

	go eventHandler.Run()

	podCollector := statistics.NewPodCollector(eventHandler)
	go podCollector.Run(clientset)

	http.Handle("/metrics", promhttp.Handler())
	handler := logging.NewHTTPHandler(http.DefaultServeMux)
	// No timeouts can be set, but that's OK for us as this HTTP server will not be
	// exposed publicly.
	//nolint:gosec
	if err := http.ListenAndServe(options.ListenAddress, handler); err != nil {
		log.Panic().Err(err).Msg(err.Error())
	}
}
