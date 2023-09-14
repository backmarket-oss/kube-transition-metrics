# Prometheus metrics

## Overview

The Prometheus metrics available expose internal details about the operation of
the `kube-transition-metrics` controller. They do not include metrics about the
pods life-cycle, those are sent as JSON data to `stdout`.

## Available metrics

An example of the available metrics can be seen below:

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
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 3.061e-05
go_gc_duration_seconds{quantile="0.25"} 4.889e-05
go_gc_duration_seconds{quantile="0.5"} 6.4192e-05
go_gc_duration_seconds{quantile="0.75"} 9.2171e-05
go_gc_duration_seconds{quantile="1"} 0.100022753
go_gc_duration_seconds_sum 24.53449887
go_gc_duration_seconds_count 2245
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 201
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.20.3"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 8.972016e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 1.4057044504e+10
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.698896e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 1.1955333e+08
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 9.19748e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 8.972016e+06
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 1.7227776e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 1.4229504e+07
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 48447
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 1.3131776e+07
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 3.145728e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.6947005743559878e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 1.19601777e+08
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 4800
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 15600
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 261120
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 408000
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 1.5300448e+07
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 866912
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 2.097152e+06
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 2.097152e+06
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 4.574132e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 8
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
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 237.28
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 10
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 5.1806208e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.6946941922e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 7.68147456e+08
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes 1.8446744073709552e+19
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 1746
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
# HELP statistic_events_handled_total Total number of statistic events handled since the last restart
# TYPE statistic_events_handled_total counter
statistic_events_handled_total 75877
```
