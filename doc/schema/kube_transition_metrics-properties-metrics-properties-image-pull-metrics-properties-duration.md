# Duration Schema

```txt
schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/image_pull/properties/duration_seconds
```

The duration in seconds to complete the image pull successfully. This is based purely off the started\_timestamp and finished\_timestamp, which themselves are based on Event timestamps which are rounded to seconds. The duration here may not match perfectly the duration seen in the kubelet image pull message, due to slight latency in reporting of image pull Events and truncation of timestamps to seconds.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [kube\_transition\_metrics.schema.json\*](kube_transition_metrics.schema.json "open original schema") |

## duration\_seconds Type

`number` ([Duration](kube_transition_metrics-properties-metrics-properties-image-pull-metrics-properties-duration.md))
