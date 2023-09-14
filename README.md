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

## Available metrics

Read about the available metrics here:
[internal/statistics/README.md](internal/statistics/README.md).

## Contributing

Send a pull request!

## Debugging

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
