# Started Timestamp Schema

```txt
schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/ready_timestamp
```

The timestamp for when the container first ready state (readinessProbe success). In the event of a pod restart, this timestamp is NOT updated.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [kube\_transition\_metrics.schema.json\*](kube_transition_metrics.schema.json "open original schema") |

## ready\_timestamp Type

`string` ([Started Timestamp](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-started-timestamp-1.md))

## ready\_timestamp Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")
