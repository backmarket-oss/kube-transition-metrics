# kube-transition-metrics architecture

## Overview

This document describes the architecture of the kube-transition-metrics
controller.
The controller watches the Kubernetes Pods API to track the life-cycle of Pods.
Pods are only tracked if they are created _after_ the startup of the controller.
This is required as certain metrics are not computable from a point-in-time
snapshot of the Kubernetes API, some timestamps are lost or replaced during the
Pod life-cycle.
Additionally, the Events API is also watched for ImagePulling and ImagePulled
events to produce the `image_pull` metric type.

## Block diagram

```mermaid
---
title: "kube-transition-metrics architecture"
---
flowchart TD
    main["./cmd/kube-transition-metrics.main()"]
    PodCollector["./internal/statistics.PodCollector.Run()"]
    StatisticEventHandler["./internal/statistics.StatisticEventHandler.Run()"]
    ImagePullCollector["./internal/statistics.ImagePullCollector.Run()"]
    Stderr["/dev/stderr"]
    HTTPServer["http.HTTPServer.ListenAndServe()"]
    PodsWatch["k8s.io/api/core/v1.PodInterface.Watch()"]
    EventsWatch["k8s.io/api/core/v1.EventInterface.Watch()"]

    main --->|"go func()"| StatisticEventHandler
    StatisticEventHandler -->|"github.com/rs/zerolog.Logger.Print()"| Stderr
    StatisticEventHandler -->|"[]./internal.statistics.podStatistic{}"| StatisticEventHandler
    PodsWatch -->|"*k8s.io/api/core/v1.Pod"| PodCollector
    main --->|"go func()"| PodCollector
    PodCollector --->|"Publish()"| StatisticEventHandler
    PodCollector --->|"go func()"| ImagePullCollector
    EventsWatch -->|"*k8s.io/api/core/v1.Event"| ImagePullCollector
    ImagePullCollector --->|"Publish()"| StatisticEventHandler
    main --> HTTPServer
```

## Goroutines

### StatisticEventHandler loop

The StatisticEventHandler goroutine reads `statisticEvent`s from a channel
written to by the PodCollector and ImagePullCollector routines.
The StatisticEventHandler manages the tracked Pods and handles `statisticEvent`s
to update these Pods in the order published.
Using a single goroutine to update the Pods statistics simplifies concurrency
control.
After processing a statistic event, it emits the statistics for the tracked pod
to stderr.

### PodCollector loop

The PodCollector goroutine receives added, modified, and deleted Pod events from
the Kubernetes API.
When Pods are added, the PodCollector sends an event to the
StatisticEventHandler to create a new tracked Pod statistic, and launches a new
ImagePullCollector routine to track Events involving the Pod UID.
When Pods are modified, the StatisticEventHandler receives an event to update
the Pod statistic.
When Pods are deleted , the StatisticEventHandler receives an event to remove
Pod statistic from tracking, and to stop the ImagePullCollector routine for this
Pod.

### ImagePullCollector loop(s)

One ImagePullCollector loop is launched by the PodCollector for each tracked
Pod.
It receives events from the Kubernetes API with the `involvedObject.uid` field
selector for the tracked Pod.
It only processes ImagePulling and ImagePulled events, and tracks the creation
timestamps of these events.
