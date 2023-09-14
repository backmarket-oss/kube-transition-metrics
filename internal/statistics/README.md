# Statistics generation

## Overview

The Pod life-cycle statistics are collection, processing, and exported as JSON
data to `stderr`.
These logs can be parsed by a log processing pipeline such as ELK or DataDog.
Debug and informative logs are mixed with metrics, and can be differentiated by
the presence of a top-level `kube_transition_metric_type`.

## Structure

All logs include a top-level `level` key, one of: `debug`, `info`, `warn`,
`error`, or `panic`.

Logs with the top-level `kube_transition_metric_type` key contain metrics about
pod life-cycle.
The `kube_transition_metrics` key contains a dictionary of metrics in such logs.

There are 3 `kube_transition_metric_type`s:

* `pod`: metrics about pod life-cycle.
* `container`: metrics about container life-cycle.
* `image_pull`: metrics about Docker/OCI image pulls.

All metric logs additionally have the `kube_namespace` and `pod_name` keys at
the top level, with the respective Kubernetes namespace and Pod name assocaited
with the metric.
The `container` and `image_pull` metric logs will additionally have the
`container_name` key set at the top level.

Metric values are only set once, and are nonexistent before the information is
available.

The metrics available for each type are as follows:

### `pod` metrics

* `scheduled_latency`: a floating-point value representing the duration in
  seconds between the Pod creation and the `PodScheduled` Condition (when a pod
  is assigned to a node).
* `initialized_latency`: a floating-point value representing the duration in
  seconds between the Pod creation and the `PodInitialized` Condition (when all
  initContainers have finished and all images are pulled).
  As a Pod can be forced to restart and initContainers are re-run, this metrics
  only represents the time to the first `PodInitialized` condition.
* `ready_latency`: a floating-point value representing the duration in seconds
  between the Pod creation and the `PodReady` Condition (when all the containers
  have started and their readinessProbe have succeeded).
  As a Pod can become unhealthy after it's first Ready, this metric only
  represents the time to the first `PodReady` Condition.

#### Example

```json
{
  "kube_transition_metrics": {
    "ready_latency": 60,
    "scheduled_latency": 1,
    "initialized_latency": 9
  },
  "level": "info",
  "kube_transition_metric_type": "pod",
  "time": "2023-09-14T14:53:50Z",
  "kube_namespace": "default",
  "pod_name": "example-pod-748d867d77-rjdd2"
}
```

### `container` metrics

* `init_container`: a boolean value representing if the container in question is
  an initContainer.
* `started_latency`: a floating-point value representing the duration in
  seconds between the Pod creation and container "start" container status (when
  the startProbe passes if it is set).
* `ready_latency`: a floating-point value representing the duration in
  seconds between the Pod creation and the "ready" container status (when the
  readinessProbe has passed).
* `running_latency`: a floating-point value representing the duration in seconds
  between the Pod creation and the "Running" container status (when a
  container's entry point is first `execve(2)`ed).

#### Example

Note: as it's relatively rare for a container to have a separate readinessProbe
and startProbe, it's not uncommon for some of these values to equal each-other.
```json
{
  "kube_transition_metrics": {
    "started_latency": 136.075222622,
    "init_container": false,
    "ready_latency": 136.075222622,
    "running_latency": 70.558999639
  },
  "container_name": "is-a-containers-container",
  "level": "info",
  "kube_transition_metric_type": "container",
  "time": "2023-09-14T15:06:39Z",
  "kube_namespace": "default",
  "pod_name": "other-example-76587b845b-xvbpn"
}
```

### `image_pull` metrics

Note: Images are only pulled when the image is not present on the node or
the imagePullPolicy is set to Always.
Images are never pulled twice for the same pod if the same image is used in
multiple different containers.
This metric will only be emitted if an ImagePulling event is initiated by the
Kubelet.

* `init_container`: a boolean value representing if the container for which the
  image was pulled is an initContainer.
* `image_pull_duration`: a floating-point value representing the duration in
  seconds between the ImagePulling event and the ImagePulled event.
  While the event message for the ImagePulled event includes a more precise but
  human-readable duration for the image pull, it is not parsed and this metric
  may be inaccurate by a few seconds as it's based on Event timestamps.

#### Example

```json
{
  "kube_transition_metrics": {
    "image_pull_duration": 16,
    "init_container": true
  },
  "container_name": "an-init-container",
  "level": "info",
  "kube_transition_metric_type": "image_pull",
  "time": "2023-09-14T15:13:22Z",
  "kube_namespace": "default",
  "pod_name": "my-pod-28245073-nzpvm"
  "message": "Successfully pulled image "ghcr.io/backmarket-oss/kube-transition-metrics:latest" in 15094.782505ms (15273.331032ms including waiting)"
}
```
