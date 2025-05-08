# Prometheus metrics

## Overview

This module offers Prometheus metrics that provide insights into the
`kube-transition-metrics` controller's internal operations.
It doesn't include pod life-cycle metrics; those are sent as JSON data to
`stdout`.

## Available metrics

Along with standard metrics from `promhttp` and `net/http/pprof`, you can see
examples of custom metrics below:

```
# HELP image_pull_collector_errors_total Total number of image pull collector errors since the last restart
# TYPE image_pull_collector_errors_total counter
image_pull_collector_errors_total 0
# HELP image_pull_collector_restarts_total Total number of times the image pull collector Watch was restarted since the process started
# TYPE image_pull_collector_restarts_total counter
image_pull_collector_restarts_total 96
# HELP image_pull_collector_routines Current number of running image pull collector routines
# TYPE image_pull_collector_routines gauge
image_pull_collector_routines 22
# HELP image_pull_statistics_tracked Current number of image pulls tracked
# TYPE image_pull_statistics_tracked gauge
image_pull_statistics_tracked 19
# HELP image_pull_watch_events_total Total number of Event Watch messages since the last restart
# TYPE image_pull_watch_events_total counter
image_pull_watch_events_total{event_type="ADDED"} 2252
image_pull_watch_events_total{event_type="MODIFIED"} 9
# HELP pod_collector_errors_total Total number of pod collector errors since the last restart
# TYPE pod_collector_errors_total counter
pod_collector_errors_total 0
# HELP pod_collector_restarts_total Total number of times the pod collector Watch was restarted since the process started
# TYPE pod_collector_restarts_total counter
pod_collector_restarts_total 0
# HELP pod_statistics_tracked Current number of pods tracked
# TYPE pod_statistics_tracked gauge
pod_statistics_tracked 114
# HELP pod_watch_events_total Total number of Pod Watch messages since the last restart
# TYPE pod_watch_events_total counter
pod_watch_events_total{event_type="ADDED"} 494
pod_watch_events_total{event_type="DELETED"} 631
pod_watch_events_total{event_type="MODIFIED"} 2903
# HELP statistic_event_processing_seconds Time spent processing events in seconds (quarantiles over 10m0s)
# TYPE statistic_event_processing_seconds summary
statistic_event_processing_seconds{event_loop="image_pull",quantile="0.5"} 0.000206833
statistic_event_processing_seconds{event_loop="image_pull",quantile="0.9"} 0.000613
statistic_event_processing_seconds{event_loop="image_pull",quantile="0.99"} 0.001455083
statistic_event_processing_seconds_sum{event_loop="image_pull"} 0.5102377143
statistic_event_processing_seconds_count{event_loop="image_pull"} 1275
statistic_event_processing_seconds{event_loop="pod",quantile="0.5"} 0.000424041
statistic_event_processing_seconds{event_loop="pod",quantile="0.9"} 0.00112025
statistic_event_processing_seconds{event_loop="pod",quantile="0.99"} 0.002443083
statistic_event_processing_seconds_sum{event_loop="pod"} 2.7123105827
statistic_event_processing_seconds_count{event_loop="pod"} 3724
# HELP statistic_event_publish_seconds Time spent waiting to publish an event in seconds (quarantiles over 10m0s)
# TYPE statistic_event_publish_seconds summary
statistic_event_publish_seconds{event_loop="image_pull",quantile="0.5"} 2.6583e-05
statistic_event_publish_seconds{event_loop="image_pull",quantile="0.9"} 4.7625e-05
statistic_event_publish_seconds{event_loop="image_pull",quantile="0.99"} 6.9e-05
statistic_event_publish_seconds_sum{event_loop="image_pull"} 0.03479350727
statistic_event_publish_seconds_count{event_loop="image_pull"} 1275
statistic_event_publish_seconds{event_loop="pod",quantile="0.5"} 1.3875e-05
statistic_event_publish_seconds{event_loop="pod",quantile="0.9"} 3.85e-05
statistic_event_publish_seconds{event_loop="pod",quantile="0.99"} 6.15e-05
statistic_event_publish_seconds_sum{event_loop="pod"} 0.06821353773
statistic_event_publish_seconds_count{event_loop="pod"} 3724
# HELP statistic_event_queue_depth Current queue depth of the event queue
# TYPE statistic_event_queue_depth gauge
statistic_event_queue_depth{event_loop="image_pull"} 0
statistic_event_queue_depth{event_loop="pod"} 0
```
