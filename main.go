package main

import (
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, _ := clientcmd.
		BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	clientset, _ := kubernetes.NewForConfig(config)

	event_handler := NewStatisticEventHandler()
	initial_sync_blacklist, resource_version, err := CollectInitialPods(clientset)
	if err != nil {
		panic(err)
	}

	event_handler.blacklistUIDs = initial_sync_blacklist
	go event_handler.Run()

	pod_collector := &PodCollector{
		eh: event_handler,
	}
	pod_collector.Run(clientset, resource_version)
}
