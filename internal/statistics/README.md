# Statistics generation

## Overview

This module handles the collection, processing, and exportation of Pod
life-cycle statistics in JSON format to `stdout`.
These logs can be integrated with log processing pipelines like ELK or DataDog.
Distinguish metric logs from debug and informative logs by the presence of a
top-level `kube_transition_metric_type`.

## Log Structure

All non-metric logs include a top-level `level` key, which can be: `debug`,
`info`, `warn`, `error`, or `panic`.
Non-metric logs are sent to `stderr`, whereas life-cycle metrics are sent to
`stdout`.

Logs with the `kube_transition_metric_type` key contain metrics about the pod
life-cycle.
The key `kube_transition_metrics` contains a dictionary of metrics within these
logs.

There are three `kube_transition_metric_type`s:
- `pod`: Metrics about pod life-cycle.
- `container`: Metrics about container life-cycle.
- `image_pull`: Metrics related to Docker/OCI image pulls.

All metric logs have the `kube_namespace` and `pod_name` keys at the top level,
which represent the Kubernetes namespace and Pod name associated with the
metric.
The `container` and `image_pull` metric logs also include the `container_name`
key.

Metrics are set once and will not be available until the information is
recorded.

The available metrics for each type are as follows:

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

Note: Images are pulled only if the image is absent on the node, or if the
`imagePullPolicy` is set to "Always".
An image won't be pulled more than once for the same pod, even if the same image
is used across multiple containers.
This metric is emitted solely when the Kubelet initiates an "ImagePulling"
event.

* `init_container`: A boolean indicating whether the container, for which the
  image was pulled, is an initContainer.
* `image_pull_duration`: Represents the duration (in seconds) between the
  "ImagePulling" and "ImagePulled" events.
  Though the "ImagePulled" event message contains a more precise, human-readable
  duration of the image pull, it isn't parsed.
  Thus, this metric's accuracy might differ slightly since it relies on event
  timestamps.

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
  "message": "Successfully pulled image \"ghcr.io/backmarket-oss/kube-transition-metrics:latest\" in 15094.782505ms (15273.331032ms including waiting)"
}
```
