# Image Pull Metrics Schema

```txt
schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/image_pull
```

Included if kube\_transition\_metric\_type is equal to "image\_pull". Note that these metrics are only emitted in the event that an image pull occurs, if imagePullPolicy is set to IfNotPresent this will only occur if the image is not already present on the node.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [kube\_transition\_metrics.schema.json\*](kube_transition_metrics.schema.json "open original schema") |

## image\_pull Type

`object` ([Image Pull Metrics](kube_transition_metrics-properties-metrics-properties-image-pull-metrics.md))

# image\_pull Properties

| Property                                   | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                                                       |
| :----------------------------------------- | :-------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [container\_name](#container_name)         | `string`  | Required | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-container-name.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/image_pull/properties/container_name")         |
| [already\_present](#already_present)       | `boolean` | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-already-present.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/image_pull/properties/already_present")       |
| [started\_timestamp](#started_timestamp)   | `string`  | Required | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-started-timestamp.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/image_pull/properties/started_timestamp")   |
| [finished\_timestamp](#finished_timestamp) | `string`  | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-finished-timestamp.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/image_pull/properties/finished_timestamp") |
| [duration\_seconds](#duration_seconds)     | `number`  | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-duration.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/image_pull/properties/duration_seconds")             |

## container\_name

The name of the container which initiated the image pull

`container_name`

* is required

* Type: `string` ([Container name](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-container-name.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-container-name.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/image_pull/properties/container_name")

### container\_name Type

`string` ([Container name](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-container-name.md))

## already\_present

true if the image was already present on the machine, otherwise false.

`already_present`

* is optional

* Type: `boolean` ([Already Present](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-already-present.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-already-present.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/image_pull/properties/already_present")

### already\_present Type

`boolean` ([Already Present](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-already-present.md))

## started\_timestamp

The timestamp for when the image pull was first initiated. This is obtained from the Event emitted by the Kubelet and may not be 100% accurate. In the event of ImagePullFailed this time is not reset for subsequent attempts.

`started_timestamp`

* is required

* Type: `string` ([Started Timestamp](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-started-timestamp.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-started-timestamp.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/image_pull/properties/started_timestamp")

### started\_timestamp Type

`string` ([Started Timestamp](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-started-timestamp.md))

### started\_timestamp Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")

## finished\_timestamp

The timestamp for when the image pull was finished. This is obtained from the Event emitted by the Kubelet and may not be 100% accurate.

`finished_timestamp`

* is optional

* Type: `string` ([Finished Timestamp](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-finished-timestamp.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-finished-timestamp.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/image_pull/properties/finished_timestamp")

### finished\_timestamp Type

`string` ([Finished Timestamp](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-finished-timestamp.md))

### finished\_timestamp Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")

## duration\_seconds

The duration in seconds to complete the image pull successfully. This is based purely off the started\_timestamp and finished\_timestamp, which themselves are based on Event timestamps which are rounded to seconds. The duration here may not match perfectly the duration seen in the kubelet image pull message, due to slight latency in reporting of image pull Events and truncation of timestamps to seconds.

`duration_seconds`

* is optional

* Type: `number` ([Duration](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-duration.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-duration.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/image_pull/properties/duration_seconds")

### duration\_seconds Type

`number` ([Duration](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-duration.md))
