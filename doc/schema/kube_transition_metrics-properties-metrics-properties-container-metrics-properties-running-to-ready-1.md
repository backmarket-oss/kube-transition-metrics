# Running to Ready Schema

```txt
schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/container/properties/started_to_ready_seconds
```

The time in seconds from the container becoming started to this container ready. Only set for non-init containers.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [kube\_transition\_metrics.schema.json\*](kube_transition_metrics.schema.json "open original schema") |

## started\_to\_ready\_seconds Type

`number` ([Running to Ready](kube_transition_metrics-properties-metrics-properties-container-metrics-properties-running-to-ready-1.md))
