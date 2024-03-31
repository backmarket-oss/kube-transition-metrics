# Pod Initialized to Running Schema

```txt
schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/initialized_to_running_seconds
```

The time in seconds from the Pod becoming initialized (all init containers exited 0) to this container running. Only set for non-init containers.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [kube\_transition\_metrics.schema.json\*](kube_transition_metrics.schema.json "open original schema") |

## initialized\_to\_running\_seconds Type

`number` ([Pod Initialized to Running](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-pod-initialized-to-running.md))
