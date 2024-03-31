# Pod Metrics Schema

```txt
schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod
```

Included if kube\_transition\_metric\_type is equal to "pod".

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [kube\_transition\_metrics.schema.json\*](kube_transition_metrics.schema.json "open original schema") |

## pod Type

`object` ([Pod Metrics](kube_transition_metrics-properties-metrics-properties-pod-metrics.md))

# pod Properties

| Property                                                                 | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                                 |
| :----------------------------------------------------------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [creation\_timestamp](#creation_timestamp)                               | `string` | Required | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-running-timestamp.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/creation_timestamp")                          |
| [scheduled\_timestamp](#scheduled_timestamp)                             | `string` | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-scheduled-timestamp.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/scheduled_timestamp")                       |
| [creation\_to\_scheduled\_seconds](#creation_to_scheduled_seconds)       | `number` | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-creation-to-scheduled.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/creation_to_scheduled_seconds")       |
| [initialized\_timestamp](#initialized_timestamp)                         | `string` | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-running-timestamp-1.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/initialized_timestamp")                     |
| [creation\_to\_initialized\_seconds](#creation_to_initialized_seconds)   | `number` | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-creation-to-initialized.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/creation_to_initialized_seconds")   |
| [scheduled\_to\_initialized\_seconds](#scheduled_to_initialized_seconds) | `number` | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-scheduled-to-initialized.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/scheduled_to_initialized_seconds") |
| [ready\_timestamp](#ready_timestamp)                                     | `string` | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-ready-timestamp.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/ready_timestamp")                               |
| [creation\_to\_ready\_seconds](#creation_to_ready_seconds)               | `number` | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-creation-to-ready.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/creation_to_ready_seconds")               |
| [initialized\_to\_ready\_seconds](#initialized_to_ready_seconds)         | `number` | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-initializing-to-running.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/initialized_to_ready_seconds")      |

## creation\_timestamp

The timestamp for when the Pod was created.

`creation_timestamp`

* is required

* Type: `string` ([Running Timestamp](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-running-timestamp.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-running-timestamp.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/creation_timestamp")

### creation\_timestamp Type

`string` ([Running Timestamp](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-running-timestamp.md))

### creation\_timestamp Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")

## scheduled\_timestamp

The timestamp for when the Pod was scheduled (Pending->Initializing state).

`scheduled_timestamp`

* is optional

* Type: `string` ([Scheduled Timestamp](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-scheduled-timestamp.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-scheduled-timestamp.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/scheduled_timestamp")

### scheduled\_timestamp Type

`string` ([Scheduled Timestamp](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-scheduled-timestamp.md))

### scheduled\_timestamp Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")

## creation\_to\_scheduled\_seconds

The time in seconds it took to schedule the Pod.

`creation_to_scheduled_seconds`

* is optional

* Type: `number` ([Pod Creation to Scheduled](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-creation-to-scheduled.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-creation-to-scheduled.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/creation_to_scheduled_seconds")

### creation\_to\_scheduled\_seconds Type

`number` ([Pod Creation to Scheduled](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-creation-to-scheduled.md))

## initialized\_timestamp

The timestamp for when the Pod first entered Running state (all init containers exited successfuly and images are pulled). In the event of a pod restart this time is not reset.

`initialized_timestamp`

* is optional

* Type: `string` ([Running Timestamp](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-running-timestamp-1.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-running-timestamp-1.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/initialized_timestamp")

### initialized\_timestamp Type

`string` ([Running Timestamp](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-running-timestamp-1.md))

### initialized\_timestamp Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")

## creation\_to\_initialized\_seconds

The time in seconds from the pod creation to when it was initialized.

`creation_to_initialized_seconds`

* is optional

* Type: `number` ([Pod Creation to Initialized](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-creation-to-initialized.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-creation-to-initialized.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/creation_to_initialized_seconds")

### creation\_to\_initialized\_seconds Type

`number` ([Pod Creation to Initialized](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-creation-to-initialized.md))

## scheduled\_to\_initialized\_seconds

The time in seconds from the pod was scheduled to when it was initialized (Initializing->Running state).

`scheduled_to_initialized_seconds`

* is optional

* Type: `number` ([Pod Scheduled to Initialized](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-scheduled-to-initialized.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-scheduled-to-initialized.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/scheduled_to_initialized_seconds")

### scheduled\_to\_initialized\_seconds Type

`number` ([Pod Scheduled to Initialized](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-scheduled-to-initialized.md))

## ready\_timestamp

The timestamp for when the Pod first became Ready (all containers had readinessProbe success). In the event of a pod restart this time is not reset.

`ready_timestamp`

* is optional

* Type: `string` ([Ready Timestamp](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-ready-timestamp.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-ready-timestamp.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/ready_timestamp")

### ready\_timestamp Type

`string` ([Ready Timestamp](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-ready-timestamp.md))

### ready\_timestamp Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")

## creation\_to\_ready\_seconds

The time in seconds from the pod creation to becoming Ready.

`creation_to_ready_seconds`

* is optional

* Type: `number` ([Pod Creation to Ready](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-creation-to-ready.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-creation-to-ready.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/creation_to_ready_seconds")

### creation\_to\_ready\_seconds Type

`number` ([Pod Creation to Ready](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-creation-to-ready.md))

## initialized\_to\_ready\_seconds

The time in seconds from the pod was initialized (Running state) to when it first bacame Ready.

`initialized_to_ready_seconds`

* is optional

* Type: `number` ([Pod Initializing to Running](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-initializing-to-running.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-initializing-to-running.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/initialized_to_ready_seconds")

### initialized\_to\_ready\_seconds Type

`number` ([Pod Initializing to Running](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-initializing-to-running.md))
