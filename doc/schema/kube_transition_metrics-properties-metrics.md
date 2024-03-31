# Metrics Schema

```txt
schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics
```

The metrics pertaining to pod\_name

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [kube\_transition\_metrics.schema.json\*](kube_transition_metrics.schema.json "open original schema") |

## kube\_transition\_metrics Type

`object` ([Metrics](kube_transition_metrics-properties-metrics.md))

all of

* [Untitled undefined type in Metric Record](kube_transition_metrics-properties-metrics-allof-0.md "check type definition")

* one (and only one) of

  * [Untitled undefined type in Metric Record](kube_transition_metrics-properties-metrics-allof-1-oneof-0.md "check type definition")

  * [Untitled undefined type in Metric Record](kube_transition_metrics-properties-metrics-allof-1-oneof-1.md "check type definition")

  * [Untitled undefined type in Metric Record](kube_transition_metrics-properties-metrics-allof-1-oneof-2.md "check type definition")

# kube\_transition\_metrics Properties

| Property                           | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                      |
| :--------------------------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [type](#type)                      | `string` | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-metric-type.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/type")                         |
| [kube\_namespace](#kube_namespace) | `string` | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-kubernetes-namespace-name.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/kube_namespace") |
| [pod\_name](#pod_name)             | `string` | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-kubernetes-pod-name.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod_name")             |
| [pod](#pod)                        | `object` | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod")                          |
| [container](#container)            | `object` | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container")              |
| [image\_pull](#image_pull)         | `object` | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-image-pull-metrics.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/image_pull")            |

## type

The type of metric included in kube\_transition\_metrics

`type`

* is optional

* Type: `string` ([Metric type](kube_transition_metrics-properties-metrics-properties-metric-type.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-metric-type.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/type")

### type Type

`string` ([Metric type](kube_transition_metrics-properties-metrics-properties-metric-type.md))

### type Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value          | Explanation |
| :------------- | :---------- |
| `"pod"`        |             |
| `"container"`  |             |
| `"image_pull"` |             |

## kube\_namespace

The name of the Kubernetes Namespace containing the pod

`kube_namespace`

* is optional

* Type: `string` ([Kubernetes Namespace name](kube_transition_metrics-properties-metrics-properties-kubernetes-namespace-name.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-kubernetes-namespace-name.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/kube_namespace")

### kube\_namespace Type

`string` ([Kubernetes Namespace name](kube_transition_metrics-properties-metrics-properties-kubernetes-namespace-name.md))

## pod\_name

The name of the Kubernetes Pod to which metrics pertain

`pod_name`

* is optional

* Type: `string` ([Kubernetes Pod name](kube_transition_metrics-properties-metrics-properties-kubernetes-pod-name.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-kubernetes-pod-name.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod_name")

### pod\_name Type

`string` ([Kubernetes Pod name](kube_transition_metrics-properties-metrics-properties-kubernetes-pod-name.md))

## pod

Included if kube\_transition\_metric\_type is equal to "pod".

`pod`

* is optional

* Type: `object` ([Pod Metrics](kube_transition_metrics-properties-metrics-properties-pod-metrics.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod")

### pod Type

`object` ([Pod Metrics](kube_transition_metrics-properties-metrics-properties-pod-metrics.md))

## container

Included if kube\_transition\_metric\_type is equal to "container".

`container`

* is optional

* Type: `object` ([Container Metrics](kube_transition_metrics-properties-metrics-properties-container-metrics.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container")

### container Type

`object` ([Container Metrics](kube_transition_metrics-properties-metrics-properties-container-metrics.md))

## image\_pull

Included if kube\_transition\_metric\_type is equal to "image\_pull". Note that these metrics are only emitted in the event that an image pull occurs, if imagePullPolicy is set to IfNotPresent this will only occur if the image is not already present on the node.

`image_pull`

* is optional

* Type: `object` ([Image Pull Metrics](kube_transition_metrics-properties-metrics-properties-image-pull-metrics.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-image-pull-metrics.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/image_pull")

### image\_pull Type

`object` ([Image Pull Metrics](kube_transition_metrics-properties-metrics-properties-image-pull-metrics.md))
