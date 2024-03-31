# Pod Scheduled to Initialized Schema

```txt
schemas/kube_transition_metrics.schema.json#/properties/kube_transition_metrics/properties/pod/properties/scheduled_to_initialized_seconds
```

The time in seconds from the pod was scheduled to when it was initialized (Initializing->Running state).

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [kube\_transition\_metrics.schema.json\*](kube_transition_metrics.schema.json "open original schema") |

## scheduled\_to\_initialized\_seconds Type

`number` ([Pod Scheduled to Initialized](kube_transition_metrics-properties-metrics-properties-pod-metrics-properties-pod-scheduled-to-initialized.md))
