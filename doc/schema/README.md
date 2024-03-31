# README

## Top-level Schemas

* [Metric Record](./kube_transition_metrics.md "JSON schema for metric logs emitted by the kube-transition-metrics controller") – `schemas/kube_transition_metrics.schema.json`

## Other Schemas

### Objects

* [Container Metrics](./kube_transition_metrics-properties-metrics-properties-container-metrics.md "Included if kube_transition_metric_type is equal to \"container\"") – `schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container`

* [Image Pull Metrics](./kube_transition_metrics-properties-metrics-properties-image-pull-metrics.md "Included if kube_transition_metric_type is equal to \"image_pull\"") – `schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/image_pull`

* [Metrics](./kube_transition_metrics-properties-metrics.md "The metrics pertaining to pod_name") – `schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics`

* [Pod Metrics](./kube_transition_metrics-properties-metrics-properties-pod-metrics.md "Included if kube_transition_metric_type is equal to \"pod\"") – `schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod`

### Arrays



## Version Note

The schemas linked above follow the JSON Schema Spec version: `https://json-schema.org/draft/2020-12/schema`
