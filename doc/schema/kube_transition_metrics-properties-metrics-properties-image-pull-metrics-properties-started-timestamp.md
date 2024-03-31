# Started Timestamp Schema

```txt
schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/image_pull/properties/started_timestamp
```

The timestamp for when the image pull was first initiated. This is obtained from the Event emitted by the Kubelet and may not be 100% accurate. In the event of ImagePullFailed this time is not reset for subsequent attempts.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [kube\_transition\_metrics.schema.json\*](kube_transition_metrics.schema.json "open original schema") |

## started\_timestamp Type

`string` ([Started Timestamp](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-started-timestamp.md))

## started\_timestamp Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")
