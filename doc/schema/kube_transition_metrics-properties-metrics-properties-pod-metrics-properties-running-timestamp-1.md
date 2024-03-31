# Running Timestamp Schema

```txt
schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/initialized_timestamp
```

The timestamp for when the Pod first entered Running state (all init containers exited successfuly and images are pulled). In the event of a pod restart this time is not reset.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [kube\_transition\_metrics.schema.json\*](kube_transition_metrics.schema.json "open original schema") |

## initialized\_timestamp Type

`string` ([Running Timestamp](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-running-timestamp-1.md))

## initialized\_timestamp Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")
