# Running to Started Schema

```txt
schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/running_to_started_seconds
```

The time in seconds from the container becoming running to this container started. Only set for non-init containers.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [kube\_transition\_metrics.schema.json\*](kube_transition_metrics.schema.json "open original schema") |

## running\_to\_started\_seconds Type

`number` ([Running to Started](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-running-to-started.md))
