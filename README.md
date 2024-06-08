# kube-transition-metrics

(WIP) Real-time statistics on the pod life-cycle timeline.

## Description

A Kubernetes controller that emits JSON logs with detailed statistics on pod
transition durations—from pod creation to readiness.
It also includes image pull statistics.

**⚠️ NOTE: This project is still in development and not ready for
production.**

## Getting Started

To use this, you need a Kubernetes cluster.
For local testing, use KIND or set it up on a remote cluster.
By default, the controller uses the current context in your `~/.kube/config`
file (i.e., the cluster shown by `kubectl cluster-info`).

```sh
go run .
```

Refer to
[cmd/kube-transition-metrics/README.md](cmd/kube-transition-metrics/README.md)
for detailed usage instructions.

To deploy this controller in-cluster, apply the helm chart located in
`charts/kube-transition-metrics`.

```sh
helm install \
    --values ./charts/kube-transition-metrics/values.yaml \
    --namespace kube-monitoring \
    kube-transition-metrics ./charts/kube-transition-metrics
```

For setting resource requests and limits, consider:

```yaml
resources:
  limits:
    cpu: 100m
    memory: 64Mi
  requests:
    cpu: 20m
    memory: 58Mi
```

To configure the DataDog agent to scrape Prometheus/OpenMetrics metrics, use
these annotations:

```yaml
podAnnotations:
  ad.datadoghq.com/kube-transition-metrics.checks: |
    {
      "openmetrics": {
        "init_config": {},
        "instances": [{
          "openmetrics_endpoint": "http://%%host%%:8080/metrics",
          "namespace": "kube_transition",
          "metrics": [".*"]
        }]
      }
    }
```

## Available metrics

This Pod life-cycle statistics are emitted in JSON format to `stdout`.
These logs can be integrated with log processing pipelines like ELK or DataDog.
Distinguish metric logs from debug and informative logs by the presence of a
top-level `kube_transition_metrics` key.

### Log Structure

All non-metric logs include a top-level `level` key, which can be: `debug`,
`info`, `warn`, `error`, or `panic`.
Non-metric logs are sent to `stderr`, whereas life-cycle metrics are sent to
`stdout`.

### Examples:

A complete pod record:

```json
{
  "kube_transition_metrics": {
    "type": "pod",
    "kube_namespace": "default",
    "pod_name": "flat-earth",
    "pod": {
      "creation_timestamp": "2024-06-08T11:14:00+02:00",
      "scheduled_timestamp": "2024-06-08T11:14:00+02:00",
      "creation_to_scheduled_seconds": 0,
      "initialized_timestamp": "2024-06-08T11:14:01+02:00",
      "creation_to_initialized_seconds": 1,
      "scheduled_to_initialized_seconds": 1,
      "ready_timestamp": "2024-06-08T11:14:02+02:00",
      "creation_to_ready_seconds": 2,
      "initialized_to_ready_seconds": 1
    }
  },
  "time": "2024-06-08T11:14:06+02:00"
}
```

A complete non-init container record:
```json
{
  "kube_transition_metrics": {
    "type": "container",
    "kube_namespace": "default",
    "pod_name": "flat-earth",
    "container": {
      "name": "consipire",
      "init_container": false,
      "initialized_to_running_seconds": 2.785652,
      "running_timestamp": "2024-06-08T11:14:02+02:00",
      "started_timestamp": "2024-06-08T11:14:02+02:00",
      "running_to_started_seconds": 0,
      "ready_timestamp": "2024-06-08T11:14:02+02:00",
      "running_to_ready_seconds": 0,
      "started_to_ready_seconds": 0
    }
  },
  "time": "2024-06-08T11:14:06+02:00"
}
```

A complete init container record:
```json
{
  "kube_transition_metrics": {
    "type": "container",
    "kube_namespace": "default",
    "pod_name": "flat-earth",
    "container": {
      "name": "subliminal-messaging",
      "init_container": true,
      "ready_timestamp": "2024-06-08T11:14:01+02:00"
    }
  },
  "time": "2024-06-08T11:14:06+02:00"
}
```

An image pull record:
```json
{
  "kube_transition_metrics": {
    "type": "image_pull",
    "image_pull": {
      "container_name": "conspire",
      "already_present": true,
      "started_timestamp": "2024-06-08T11:14:01+02:00",
      "finished_timestamp": "2024-06-08T11:14:01+02:00",
      "duration_seconds": 0
    },
    "kube_namespace": "default",
    "pod_name": "flat-earth"
  },
  "time": "2024-06-08T11:14:01+02:00",
  "message": "Container image \"docker.io/library/nginx:latest\" already present on machine"
}
```

For a detailed overview of available metrics, see [doc/SCHEMA.md](doc/SCHEMA.md).

## Contributing

We welcome contributions! Please send a pull request.

For a comprehensive overview, read the [Architecture](doc/ARCHITECTURE.md)
design document.

Internal metrics about the controller's operations are published using
`promhttp`.
These metrics are separate from the pod life-cycle statistics.
For more information on internal observability, see
[internal/prommetrics/README.md](internal/prommetrics/README.md).

Profiling is also enabled through the `/debug/pprof/` endpoints.
Refer to [net/http/pprof](https://pkg.go.dev/net/http/pprof).

## License
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License"); you may not usex
this file except in compliance with the License.
You may obtain a copy of the License at
[https://www.apache.org/licenses/LICENSE-2.0](https://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
