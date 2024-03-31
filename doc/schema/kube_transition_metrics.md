# Metric Record Schema

```txt
schemas/kube_transition_metrics.schema.json
```

JSON schema for metric logs emitted by the kube-transition-metrics controller

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [kube\_transition\_metrics.schema.json](kube_transition_metrics.schema.json "open original schema") |

## Metric Record Type

`object` ([Metric Record](kube_transition_metrics.md))

# Metric Record Properties

| Property                                              | Type     | Required | Nullable       | Defined by                                                                                                                                       |
| :---------------------------------------------------- | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------- |
| [kube\_transition\_metrics](#kube_transition_metrics) | Merged   | Required | cannot be null | [Metric Record](kube_transition_metrics-properties-metrics.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics") |
| [time](#time)                                         | `string` | Required | cannot be null | [Metric Record](kube_transition_metrics-properties-metric-timestamp.md "schemas/kube_transition_metrics.schema.json#/properties/time")           |
| [message](#message)                                   | `string` | Optional | cannot be null | [Metric Record](kube_transition_metrics-properties-message.md "schemas/kube_transition_metrics.schema.json#/properties/message")                 |

## kube\_transition\_metrics

The metrics pertaining to pod\_name

`kube_transition_metrics`

* is required

* Type: `object` ([Metrics](kube_transition_metrics-properties-metrics.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metrics.md "schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics")

### kube\_transition\_metrics Type

`object` ([Metrics](kube_transition_metrics-properties-metrics.md))

all of

* [Untitled undefined type in Metric Record](kube_transition_metrics-properties-metrics-allof-0.md "check type definition")

* one (and only one) of

  * [Untitled undefined type in Metric Record](kube_transition_metrics-properties-metrics-allof-1-oneof-0.md "check type definition")

  * [Untitled undefined type in Metric Record](kube_transition_metrics-properties-metrics-allof-1-oneof-1.md "check type definition")

  * [Untitled undefined type in Metric Record](kube_transition_metrics-properties-metrics-allof-1-oneof-2.md "check type definition")

## time

The time at which this metric was emitted.

`time`

* is required

* Type: `string` ([Metric Timestamp](kube_transition_metrics-properties-metric-timestamp.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-metric-timestamp.md "schemas/kube_transition_metrics.schema.json#/properties/time")

### time Type

`string` ([Metric Timestamp](kube_transition_metrics-properties-metric-timestamp.md))

### time Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")

## message

An additional message emitted along with metrics.

`message`

* is optional

* Type: `string` ([Message](kube_transition_metrics-properties-message.md))

* cannot be null

* defined in: [Metric Record](kube_transition_metrics-properties-message.md "schemas/kube_transition_metrics.schema.json#/properties/message")

### message Type

`string` ([Message](kube_transition_metrics-properties-message.md))
