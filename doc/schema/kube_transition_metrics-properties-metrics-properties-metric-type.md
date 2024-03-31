# Metric type Schema

```txt
schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/type
```

The type of metric included in kube\_transition\_metrics

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [kube\_transition\_metrics.schema.json\*](kube_transition_metrics.schema.json "open original schema") |

## type Type

`string` ([Metric type](kube_transition_metrics-properties-metrics-properties-metric-type.md))

## type Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value          | Explanation |
| :------------- | :---------- |
| `"pod"`        |             |
| `"container"`  |             |
| `"image_pull"` |             |
