# Pod Initializing to Running Schema

```txt
schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/initialized_to_ready_seconds
```

The time in seconds from the pod was initialized (Running state) to when it first bacame Ready.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [kube\_transition\_metrics.schema.json\*](kube_transition_metrics.schema.json "open original schema") |

## initialized\_to\_ready\_seconds Type

`number` ([Pod Initializing to Running](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-initializing-to-running.md))
