# Previous Container Finished to Running Schema

```txt
schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/previous_to_running_seconds
```

The time in seconds from the previous init container becoming Ready (exited 0) to this container running. Only set for init containers, absent for the first init container.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [kube\_transition\_metrics.schema.json\*](kube_transition_metrics.schema.json "open original schema") |

## previous\_to\_running\_seconds Type

`number` ([Previous Container Finished to Running](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-previous-container-finished-to-running.md))
