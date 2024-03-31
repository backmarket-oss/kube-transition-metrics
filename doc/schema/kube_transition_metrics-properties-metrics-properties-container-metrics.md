# Container Metrics Schema

```txt
schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container
```

Included if kube\_transition\_metric\_type is equal to "container".

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [kube\_transition\_metrics.schema.json\*](kube_transition_metrics.schema.json "open original schema") |

## container Type

`object` ([Container Metrics](kube_transition_metrics-properties-metrics-properties-container-metrics.md))

# container Properties

| Property                                                             | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                                                                                  |
| :------------------------------------------------------------------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [name](#name)                                                        | `string`  | Required | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-container-name.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/name")                                                |
| [init\_container](#init_container)                                   | `boolean` | Required | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-init-container.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/init_container")                                      |
| [previous\_to\_running\_seconds](#previous_to_running_seconds)       | `number`  | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-previous-container-finished-to-running.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/previous_to_running_seconds") |
| [initialized\_to\_running\_seconds](#initialized_to_running_seconds) | `number`  | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-pod-initialized-to-running.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/initialized_to_running_seconds")          |
| [running\_timestamp](#running_timestamp)                             | `string`  | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-running-timestamp.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/running_timestamp")                                |
| [started\_timestamp](#started_timestamp)                             | `string`  | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-started-timestamp.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/started_timestamp")                                |
| [running\_to\_started\_seconds](#running_to_started_seconds)         | `number`  | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-running-to-started.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/running_to_started_seconds")                      |
| [ready\_timestamp](#ready_timestamp)                                 | `string`  | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-started-timestamp-1.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/ready_timestamp")                                |
| [running\_to\_ready\_seconds](#running_to_ready_seconds)             | `number`  | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-running-to-ready.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/running_to_ready_seconds")                          |
| [started\_to\_ready\_seconds](#started_to_ready_seconds)             | `number`  | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-running-to-ready-1.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/started_to_ready_seconds")                        |

## name

The name of the container to which metrics pertain

`name`

* is required

* Type: `string` ([Container name](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-container-name.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-container-name.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/name")

### name Type

`string` ([Container name](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-container-name.md))

## init\_container

True if the container is an init container, otherwise false.

`init_container`

* is required

* Type: `boolean` ([Init Container](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-init-container.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-init-container.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/init_container")

### init\_container Type

`boolean` ([Init Container](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-init-container.md))

## previous\_to\_running\_seconds

The time in seconds from the previous init container becoming Ready (exited 0) to this container running. Only set for init containers, absent for the first init container.

`previous_to_running_seconds`

* is optional

* Type: `number` ([Previous Container Finished to Running](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-previous-container-finished-to-running.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-previous-container-finished-to-running.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/previous_to_running_seconds")

### previous\_to\_running\_seconds Type

`number` ([Previous Container Finished to Running](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-previous-container-finished-to-running.md))

## initialized\_to\_running\_seconds

The time in seconds from the Pod becoming initialized (all init containers exited 0) to this container running. Only set for non-init containers.

`initialized_to_running_seconds`

* is optional

* Type: `number` ([Pod Initialized to Running](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-pod-initialized-to-running.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-pod-initialized-to-running.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/initialized_to_running_seconds")

### initialized\_to\_running\_seconds Type

`number` ([Pod Initialized to Running](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-pod-initialized-to-running.md))

## running\_timestamp

The timestamp for when the container first entered Running state (first fork(2)/execve(2) in container environment). In the event of a pod restart, this timestamp is NOT updated.

`running_timestamp`

* is optional

* Type: `string` ([Running Timestamp](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-running-timestamp.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-running-timestamp.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/running_timestamp")

### running\_timestamp Type

`string` ([Running Timestamp](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-running-timestamp.md))

### running\_timestamp Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")

## started\_timestamp

The timestamp for when the container first started state (startupProbe success). In the event of a pod restart, this timestamp is NOT updated. Only set for non-init containers.

`started_timestamp`

* is optional

* Type: `string` ([Started Timestamp](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-started-timestamp.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-started-timestamp.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/started_timestamp")

### started\_timestamp Type

`string` ([Started Timestamp](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-started-timestamp.md))

### started\_timestamp Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")

## running\_to\_started\_seconds

The time in seconds from the container becoming running to this container started. Only set for non-init containers.

`running_to_started_seconds`

* is optional

* Type: `number` ([Running to Started](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-running-to-started.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-running-to-started.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/running_to_started_seconds")

### running\_to\_started\_seconds Type

`number` ([Running to Started](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-running-to-started.md))

## ready\_timestamp

The timestamp for when the container first ready state (readinessProbe success). In the event of a pod restart, this timestamp is NOT updated.

`ready_timestamp`

* is optional

* Type: `string` ([Started Timestamp](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-started-timestamp-1.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-started-timestamp-1.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/ready_timestamp")

### ready\_timestamp Type

`string` ([Started Timestamp](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-started-timestamp-1.md))

### ready\_timestamp Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")

## running\_to\_ready\_seconds

The time in seconds from the container becoming running to this container ready. In init containers, this is the time the container exited with a successful status.

`running_to_ready_seconds`

* is optional

* Type: `number` ([Running to Ready](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-running-to-ready.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-running-to-ready.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/running_to_ready_seconds")

### running\_to\_ready\_seconds Type

`number` ([Running to Ready](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-running-to-ready.md))

## started\_to\_ready\_seconds

The time in seconds from the container becoming started to this container ready. Only set for non-init containers.

`started_to_ready_seconds`

* is optional

* Type: `number` ([Running to Ready](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-running-to-ready-1.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-running-to-ready-1.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/started_to_ready_seconds")

### started\_to\_ready\_seconds Type

`number` ([Running to Ready](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-running-to-ready-1.md))
