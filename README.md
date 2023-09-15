# kube-transition-metrics

(WIP) Live statistics on pod life-cycle timeline.

## Description

Kubernetes controller emitting JSON logs containing granular statistics of pod
transition durations for pod creation to pod Ready.
Includes image pull statistics as well.

**WORK IN PROGRESS, NOT PRODUCTION READY**

## Getting Started

You’ll need a Kubernetes cluster to run against.
You can use KIND to get a local cluster for testing, or run against a remote
cluster. Note: The controller will automatically use the current context in
your `~/.kube/config` file (i.e. whatever cluster kubectl cluster-info shows).

```sh
go run .
```

See:
[cmd/kube-transition-metrics/README.md](cmd/kube-transition-metrics/README.md)
for more information on usage.

The helm chart in `charts/kube-transition-metrics` can be applied to run this
controller in-cluster.

```sh
helm install \
    --values ./charts/kube-transition-metrics/values.yaml \
    --namespace kube-monitoring \
    kube-transition-metrics ./charts/kube-transition-metrics
```

Reasonable resource requests and limits may be:

```yaml
resources:
  limits:
    cpu: 100m
    memory: 64Mi
  requests:
    cpu: 20m
    memory: 58Mi
```

Annotations to configure DataDog agent to scrape Prometheus/OpenMetrics metrics
are as follows:

```yaml
podAnnotations:
  ad.datadoghq.com/kube-transition-monitoring.checks: |
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

Read about the available metrics here:
[internal/statistics/README.md](internal/statistics/README.md).

## Contributing

Send a pull request!

## Hacking

Read the [Architecture](doc/ARCHITECTURE.md) design document for a high-level
overview.

`promhttp` is used to publish some metrics about the internals of the
controller.
These metrics are distinct from the pod life-cycle timeline statistics.

See details about the internal observability here:
[internal/prommetrics/README.md](internal/prommetrics/README.md).

`pprof` is also instrumented through the `/debug/pprof/` endpoints, see
[net/http/pprof](https://pkg.go.dev/net/http/pprof).

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
