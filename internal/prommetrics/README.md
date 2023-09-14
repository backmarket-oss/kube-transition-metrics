# Prometheus metrics

## Overview

The Prometheus metrics available expose internal details about the operation of
the `kube-transition-metrics` controller. They do not include metrics about the
pods life-cycle, those are sent as JSON data to `stdout`.

## Available metrics

In addition to the standard metrics instrumented out of the box with `promhttp`
and `net/http/pprof`, an example of the custom metrics can be seen below:

```
# HELP channel_monitors Current number of channel monitor goroutines
# TYPE channel_monitors gauge
channel_monitors 1
# HELP channel_publish_wait_total Total amount of time in seconds waiting to publish a to the channel
# TYPE channel_publish_wait_total counter
channel_publish_wait_total{channel_name="statistic_events"} 13.42432692799992
# HELP channel_queue_depth Current queue depth of the channel
# TYPE channel_queue_depth gauge
channel_queue_depth{channel_name="statistic_events"} 0
# HELP events_watch_processed_total Total number of Event Watch messages since the last restart
# TYPE events_watch_processed_total counter
events_watch_processed_total{event_type="ADDED"} 40736
events_watch_processed_total{event_type="DELETED"} 275
events_watch_processed_total{event_type="MODIFIED"} 2082
# HELP image_pull_collector_errors_total Total number of image pull collector errors since the last restart
# TYPE image_pull_collector_errors_total counter
image_pull_collector_errors_total 0
# HELP image_pull_collector_restarts_total Total number of times the image pull collector Watch was restarted since the process started
# TYPE image_pull_collector_restarts_total counter
image_pull_collector_restarts_total 4058
# HELP image_pull_collector_routines Current number of running image pull collector routines
# TYPE image_pull_collector_routines gauge
image_pull_collector_routines 63
# HELP pod_collector_errors_total Total number of pod collector errors since the last restart
# TYPE pod_collector_errors_total counter
pod_collector_errors_total 0
# HELP pod_collector_restarts_total Total number of times the pod collector Watch was restarted since the process started
# TYPE pod_collector_restarts_total counter
pod_collector_restarts_total 277
# HELP pod_statistics_tracked Current number of pods tracked
# TYPE pod_statistics_tracked gauge
pod_statistics_tracked 833
# HELP pods_watch_processed_total Total number of Pod Watch messages since the last restart
# TYPE pods_watch_processed_total counter
pods_watch_processed_total{event_type="ADDED"} 179643
pods_watch_processed_total{event_type="DELETED"} 1091
pods_watch_processed_total{event_type="MODIFIED"} 10532
# HELP statistic_events_handled_total Total number of statistic events handled since the last restart
# TYPE statistic_events_handled_total counter
statistic_events_handled_total 75877
```
