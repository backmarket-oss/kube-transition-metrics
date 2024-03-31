# Running to Ready Schema

```txt
schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/running_to_ready_seconds
```

The time in seconds from the container becoming running to this container ready. In init containers, this is the time the container exited with a successful status.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [kube\_transition\_metrics.schema.json\*](kube_transition_metrics.schema.json "open original schema") |

## running\_to\_ready\_seconds Type

`number` ([Running to Ready](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-running-to-ready.md))
