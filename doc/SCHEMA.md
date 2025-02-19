# Metric Record

- [1. Property `Metric Record > kube_transition_metrics`](#kube_transition_metrics)
  - [1.1. Property `Metric Record > kube_transition_metrics > allOf > item 0`](#kube_transition_metrics_allOf_i0)
    - [1.1.1. The following properties are required](#autogenerated_heading_2)
  - [1.2. Property `Metric Record > kube_transition_metrics > allOf > item 1`](#kube_transition_metrics_allOf_i1)
    - [1.2.1. Property `Metric Record > kube_transition_metrics > allOf > item 1 > oneOf > item 0`](#kube_transition_metrics_allOf_i1_oneOf_i0)
      - [1.2.1.1. The following properties are required](#autogenerated_heading_3)
    - [1.2.2. Property `Metric Record > kube_transition_metrics > allOf > item 1 > oneOf > item 1`](#kube_transition_metrics_allOf_i1_oneOf_i1)
      - [1.2.2.1. The following properties are required](#autogenerated_heading_4)
    - [1.2.3. Property `Metric Record > kube_transition_metrics > allOf > item 1 > oneOf > item 2`](#kube_transition_metrics_allOf_i1_oneOf_i2)
      - [1.2.3.1. The following properties are required](#autogenerated_heading_5)
  - [1.3. Property `Metric Record > kube_transition_metrics > type`](#kube_transition_metrics_type)
  - [1.4. Property `Metric Record > kube_transition_metrics > kube_namespace`](#kube_transition_metrics_kube_namespace)
  - [1.5. Property `Metric Record > kube_transition_metrics > pod_name`](#kube_transition_metrics_pod_name)
  - [1.6. Property `Metric Record > kube_transition_metrics > pod`](#kube_transition_metrics_pod)
    - [1.6.1. Property `Metric Record > kube_transition_metrics > pod > creation_timestamp`](#kube_transition_metrics_pod_creation_timestamp)
    - [1.6.2. Property `Metric Record > kube_transition_metrics > pod > scheduled_timestamp`](#kube_transition_metrics_pod_scheduled_timestamp)
    - [1.6.3. Property `Metric Record > kube_transition_metrics > pod > creation_to_scheduled_seconds`](#kube_transition_metrics_pod_creation_to_scheduled_seconds)
    - [1.6.4. Property `Metric Record > kube_transition_metrics > pod > initialized_timestamp`](#kube_transition_metrics_pod_initialized_timestamp)
    - [1.6.5. Property `Metric Record > kube_transition_metrics > pod > creation_to_initialized_seconds`](#kube_transition_metrics_pod_creation_to_initialized_seconds)
    - [1.6.6. Property `Metric Record > kube_transition_metrics > pod > scheduled_to_initialized_seconds`](#kube_transition_metrics_pod_scheduled_to_initialized_seconds)
    - [1.6.7. Property `Metric Record > kube_transition_metrics > pod > ready_timestamp`](#kube_transition_metrics_pod_ready_timestamp)
    - [1.6.8. Property `Metric Record > kube_transition_metrics > pod > creation_to_ready_seconds`](#kube_transition_metrics_pod_creation_to_ready_seconds)
    - [1.6.9. Property `Metric Record > kube_transition_metrics > pod > initialized_to_ready_seconds`](#kube_transition_metrics_pod_initialized_to_ready_seconds)
  - [1.7. Property `Metric Record > kube_transition_metrics > container`](#kube_transition_metrics_container)
    - [1.7.1. Property `Metric Record > kube_transition_metrics > container > name`](#kube_transition_metrics_container_name)
    - [1.7.2. Property `Metric Record > kube_transition_metrics > container > init_container`](#kube_transition_metrics_container_init_container)
    - [1.7.3. Property `Metric Record > kube_transition_metrics > container > previous_to_running_seconds`](#kube_transition_metrics_container_previous_to_running_seconds)
    - [1.7.4. Property `Metric Record > kube_transition_metrics > container > initialized_to_running_seconds`](#kube_transition_metrics_container_initialized_to_running_seconds)
    - [1.7.5. Property `Metric Record > kube_transition_metrics > container > running_timestamp`](#kube_transition_metrics_container_running_timestamp)
    - [1.7.6. Property `Metric Record > kube_transition_metrics > container > started_timestamp`](#kube_transition_metrics_container_started_timestamp)
    - [1.7.7. Property `Metric Record > kube_transition_metrics > container > running_to_started_seconds`](#kube_transition_metrics_container_running_to_started_seconds)
    - [1.7.8. Property `Metric Record > kube_transition_metrics > container > ready_timestamp`](#kube_transition_metrics_container_ready_timestamp)
    - [1.7.9. Property `Metric Record > kube_transition_metrics > container > running_to_ready_seconds`](#kube_transition_metrics_container_running_to_ready_seconds)
    - [1.7.10. Property `Metric Record > kube_transition_metrics > container > started_to_ready_seconds`](#kube_transition_metrics_container_started_to_ready_seconds)
  - [1.8. Property `Metric Record > kube_transition_metrics > image_pull`](#kube_transition_metrics_image_pull)
    - [1.8.1. Property `Metric Record > kube_transition_metrics > image_pull > container_name`](#kube_transition_metrics_image_pull_container_name)
    - [1.8.2. Property `Metric Record > kube_transition_metrics > image_pull > already_present`](#kube_transition_metrics_image_pull_already_present)
    - [1.8.3. Property `Metric Record > kube_transition_metrics > image_pull > started_timestamp`](#kube_transition_metrics_image_pull_started_timestamp)
    - [1.8.4. Property `Metric Record > kube_transition_metrics > image_pull > finished_timestamp`](#kube_transition_metrics_image_pull_finished_timestamp)
    - [1.8.5. Property `Metric Record > kube_transition_metrics > image_pull > duration_seconds`](#kube_transition_metrics_image_pull_duration_seconds)
- [2. Property `Metric Record > time`](#time)
- [3. Property `Metric Record > message`](#message)

**Title:** Metric Record

|                           |             |
| ------------------------- | ----------- |
| **Type**                  | `object`    |
| **Required**              | No          |
| **Additional properties** | Not allowed |

**Description:** JSON schema for metric logs emitted by the kube-transition-metrics controller

| Property                                               | Type        | Title/Description |
| ------------------------------------------------------ | ----------- | ----------------- |
| + [kube_transition_metrics](#kube_transition_metrics ) | Combination | Metrics           |
| + [time](#time )                                       | string      | Metric Timestamp  |
| - [message](#message )                                 | string      | Message           |

## <a name="kube_transition_metrics"></a>1. Property `Metric Record > kube_transition_metrics`

**Title:** Metrics

|                           |             |
| ------------------------- | ----------- |
| **Type**                  | `combining` |
| **Required**              | Yes         |
| **Additional properties** | Not allowed |

**Description:** The metrics pertaining to pod_name

| Property                                                     | Type             | Title/Description         |
| ------------------------------------------------------------ | ---------------- | ------------------------- |
| - [type](#kube_transition_metrics_type )                     | enum (of string) | Metric type               |
| - [kube_namespace](#kube_transition_metrics_kube_namespace ) | string           | Kubernetes Namespace name |
| - [pod_name](#kube_transition_metrics_pod_name )             | string           | Kubernetes Pod name       |
| - [pod](#kube_transition_metrics_pod )                       | object           | Pod Metrics               |
| - [container](#kube_transition_metrics_container )           | object           | Container Metrics         |
| - [image_pull](#kube_transition_metrics_image_pull )         | object           | Image Pull Metrics        |

| All of(Requirement)                         |
| ------------------------------------------- |
| [item 0](#kube_transition_metrics_allOf_i0) |
| [item 1](#kube_transition_metrics_allOf_i1) |

### <a name="kube_transition_metrics_allOf_i0"></a>1.1. Property `Metric Record > kube_transition_metrics > allOf > item 0`

|                           |                  |
| ------------------------- | ---------------- |
| **Type**                  | `object`         |
| **Required**              | No               |
| **Additional properties** | Any type allowed |

#### <a name="autogenerated_heading_2"></a>1.1.1. The following properties are required
* kube_namespace
* pod_name
* type

### <a name="kube_transition_metrics_allOf_i1"></a>1.2. Property `Metric Record > kube_transition_metrics > allOf > item 1`

|                           |                  |
| ------------------------- | ---------------- |
| **Type**                  | `combining`      |
| **Required**              | No               |
| **Additional properties** | Any type allowed |

| One of(Option)                                       |
| ---------------------------------------------------- |
| [item 0](#kube_transition_metrics_allOf_i1_oneOf_i0) |
| [item 1](#kube_transition_metrics_allOf_i1_oneOf_i1) |
| [item 2](#kube_transition_metrics_allOf_i1_oneOf_i2) |

#### <a name="kube_transition_metrics_allOf_i1_oneOf_i0"></a>1.2.1. Property `Metric Record > kube_transition_metrics > allOf > item 1 > oneOf > item 0`

|                           |                  |
| ------------------------- | ---------------- |
| **Type**                  | `object`         |
| **Required**              | No               |
| **Additional properties** | Any type allowed |

##### <a name="autogenerated_heading_3"></a>1.2.1.1. The following properties are required
* pod

#### <a name="kube_transition_metrics_allOf_i1_oneOf_i1"></a>1.2.2. Property `Metric Record > kube_transition_metrics > allOf > item 1 > oneOf > item 1`

|                           |                  |
| ------------------------- | ---------------- |
| **Type**                  | `object`         |
| **Required**              | No               |
| **Additional properties** | Any type allowed |

##### <a name="autogenerated_heading_4"></a>1.2.2.1. The following properties are required
* container

#### <a name="kube_transition_metrics_allOf_i1_oneOf_i2"></a>1.2.3. Property `Metric Record > kube_transition_metrics > allOf > item 1 > oneOf > item 2`

|                           |                  |
| ------------------------- | ---------------- |
| **Type**                  | `object`         |
| **Required**              | No               |
| **Additional properties** | Any type allowed |

##### <a name="autogenerated_heading_5"></a>1.2.3.1. The following properties are required
* image_pull

### <a name="kube_transition_metrics_type"></a>1.3. Property `Metric Record > kube_transition_metrics > type`

**Title:** Metric type

|              |                    |
| ------------ | ------------------ |
| **Type**     | `enum (of string)` |
| **Required** | No                 |

**Description:** The type of metric included in kube_transition_metrics

Must be one of:
* "pod"
* "container"
* "image_pull"

### <a name="kube_transition_metrics_kube_namespace"></a>1.4. Property `Metric Record > kube_transition_metrics > kube_namespace`

**Title:** Kubernetes Namespace name

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |

**Description:** The name of the Kubernetes Namespace containing the pod

### <a name="kube_transition_metrics_pod_name"></a>1.5. Property `Metric Record > kube_transition_metrics > pod_name`

**Title:** Kubernetes Pod name

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |

**Description:** The name of the Kubernetes Pod to which metrics pertain

### <a name="kube_transition_metrics_pod"></a>1.6. Property `Metric Record > kube_transition_metrics > pod`

**Title:** Pod Metrics

|                           |             |
| ------------------------- | ----------- |
| **Type**                  | `object`    |
| **Required**              | No          |
| **Additional properties** | Not allowed |

**Description:** Included if kube_transition_metric_type is equal to "pod".

| Property                                                                                             | Type   | Title/Description            |
| ---------------------------------------------------------------------------------------------------- | ------ | ---------------------------- |
| + [creation_timestamp](#kube_transition_metrics_pod_creation_timestamp )                             | string | Running Timestamp            |
| - [scheduled_timestamp](#kube_transition_metrics_pod_scheduled_timestamp )                           | string | Scheduled Timestamp          |
| - [creation_to_scheduled_seconds](#kube_transition_metrics_pod_creation_to_scheduled_seconds )       | number | Pod Creation to Scheduled    |
| - [initialized_timestamp](#kube_transition_metrics_pod_initialized_timestamp )                       | string | initialized Timestamp        |
| - [creation_to_initialized_seconds](#kube_transition_metrics_pod_creation_to_initialized_seconds )   | number | Pod Creation to Initialized  |
| - [scheduled_to_initialized_seconds](#kube_transition_metrics_pod_scheduled_to_initialized_seconds ) | number | Pod Scheduled to Initialized |
| - [ready_timestamp](#kube_transition_metrics_pod_ready_timestamp )                                   | string | Ready Timestamp              |
| - [creation_to_ready_seconds](#kube_transition_metrics_pod_creation_to_ready_seconds )               | number | Pod Creation to Ready        |
| - [initialized_to_ready_seconds](#kube_transition_metrics_pod_initialized_to_ready_seconds )         | number | Pod Initialized to Ready     |

#### <a name="kube_transition_metrics_pod_creation_timestamp"></a>1.6.1. Property `Metric Record > kube_transition_metrics > pod > creation_timestamp`

**Title:** Running Timestamp

|              |             |
| ------------ | ----------- |
| **Type**     | `string`    |
| **Required** | Yes         |
| **Format**   | `date-time` |

**Description:** The timestamp for when the Pod was created.

#### <a name="kube_transition_metrics_pod_scheduled_timestamp"></a>1.6.2. Property `Metric Record > kube_transition_metrics > pod > scheduled_timestamp`

**Title:** Scheduled Timestamp

|              |             |
| ------------ | ----------- |
| **Type**     | `string`    |
| **Required** | No          |
| **Format**   | `date-time` |

**Description:** The timestamp for when the Pod was scheduled (Pending->Initializing state).

#### <a name="kube_transition_metrics_pod_creation_to_scheduled_seconds"></a>1.6.3. Property `Metric Record > kube_transition_metrics > pod > creation_to_scheduled_seconds`

**Title:** Pod Creation to Scheduled

|              |          |
| ------------ | -------- |
| **Type**     | `number` |
| **Required** | No       |

**Description:** The time in seconds it took to schedule the Pod.

#### <a name="kube_transition_metrics_pod_initialized_timestamp"></a>1.6.4. Property `Metric Record > kube_transition_metrics > pod > initialized_timestamp`

**Title:** initialized Timestamp

|              |             |
| ------------ | ----------- |
| **Type**     | `string`    |
| **Required** | No          |
| **Format**   | `date-time` |

**Description:** The timestamp for when the Pod first entered Running state (all init containers exited successfuly and images are pulled). In the event of a pod restart this time is not reset.

#### <a name="kube_transition_metrics_pod_creation_to_initialized_seconds"></a>1.6.5. Property `Metric Record > kube_transition_metrics > pod > creation_to_initialized_seconds`

**Title:** Pod Creation to Initialized

|              |          |
| ------------ | -------- |
| **Type**     | `number` |
| **Required** | No       |

**Description:** The time in seconds from the pod creation to when it was initialized.

#### <a name="kube_transition_metrics_pod_scheduled_to_initialized_seconds"></a>1.6.6. Property `Metric Record > kube_transition_metrics > pod > scheduled_to_initialized_seconds`

**Title:** Pod Scheduled to Initialized

|              |          |
| ------------ | -------- |
| **Type**     | `number` |
| **Required** | No       |

**Description:** The time in seconds from the pod was scheduled to when it was initialized (Initializing->Running state).

#### <a name="kube_transition_metrics_pod_ready_timestamp"></a>1.6.7. Property `Metric Record > kube_transition_metrics > pod > ready_timestamp`

**Title:** Ready Timestamp

|              |             |
| ------------ | ----------- |
| **Type**     | `string`    |
| **Required** | No          |
| **Format**   | `date-time` |

**Description:** The timestamp for when the Pod first became Ready (all containers had readinessProbe success). In the event of a pod restart this time is not reset.

#### <a name="kube_transition_metrics_pod_creation_to_ready_seconds"></a>1.6.8. Property `Metric Record > kube_transition_metrics > pod > creation_to_ready_seconds`

**Title:** Pod Creation to Ready

|              |          |
| ------------ | -------- |
| **Type**     | `number` |
| **Required** | No       |

**Description:** The time in seconds from the pod creation to becoming Ready.

#### <a name="kube_transition_metrics_pod_initialized_to_ready_seconds"></a>1.6.9. Property `Metric Record > kube_transition_metrics > pod > initialized_to_ready_seconds`

**Title:** Pod Initialized to Ready

|              |          |
| ------------ | -------- |
| **Type**     | `number` |
| **Required** | No       |

**Description:** The time in seconds from the pod was initialized (Running state) to when it first bacame Ready.

### <a name="kube_transition_metrics_container"></a>1.7. Property `Metric Record > kube_transition_metrics > container`

**Title:** Container Metrics

|                           |             |
| ------------------------- | ----------- |
| **Type**                  | `object`    |
| **Required**              | No          |
| **Additional properties** | Not allowed |

**Description:** Included if kube_transition_metric_type is equal to "container".

| Property                                                                                               | Type    | Title/Description                      |
| ------------------------------------------------------------------------------------------------------ | ------- | -------------------------------------- |
| + [name](#kube_transition_metrics_container_name )                                                     | string  | Container name                         |
| + [init_container](#kube_transition_metrics_container_init_container )                                 | boolean | Init Container                         |
| - [previous_to_running_seconds](#kube_transition_metrics_container_previous_to_running_seconds )       | number  | Previous Container Finished to Running |
| - [initialized_to_running_seconds](#kube_transition_metrics_container_initialized_to_running_seconds ) | number  | Pod Initialized to Running             |
| - [running_timestamp](#kube_transition_metrics_container_running_timestamp )                           | string  | Running Timestamp                      |
| - [started_timestamp](#kube_transition_metrics_container_started_timestamp )                           | string  | Started Timestamp                      |
| - [running_to_started_seconds](#kube_transition_metrics_container_running_to_started_seconds )         | number  | Running to Started                     |
| - [ready_timestamp](#kube_transition_metrics_container_ready_timestamp )                               | string  | Started Timestamp                      |
| - [running_to_ready_seconds](#kube_transition_metrics_container_running_to_ready_seconds )             | number  | Running to Ready                       |
| - [started_to_ready_seconds](#kube_transition_metrics_container_started_to_ready_seconds )             | number  | Running to Ready                       |

#### <a name="kube_transition_metrics_container_name"></a>1.7.1. Property `Metric Record > kube_transition_metrics > container > name`

**Title:** Container name

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | Yes      |

**Description:** The name of the container to which metrics pertain

#### <a name="kube_transition_metrics_container_init_container"></a>1.7.2. Property `Metric Record > kube_transition_metrics > container > init_container`

**Title:** Init Container

|              |           |
| ------------ | --------- |
| **Type**     | `boolean` |
| **Required** | Yes       |

**Description:** True if the container is an init container, otherwise false.

#### <a name="kube_transition_metrics_container_previous_to_running_seconds"></a>1.7.3. Property `Metric Record > kube_transition_metrics > container > previous_to_running_seconds`

**Title:** Previous Container Finished to Running

|              |          |
| ------------ | -------- |
| **Type**     | `number` |
| **Required** | No       |

**Description:** The time in seconds from the previous init container becoming Ready (exited 0) to this container running. Only set for init containers, absent for the first init container.

#### <a name="kube_transition_metrics_container_initialized_to_running_seconds"></a>1.7.4. Property `Metric Record > kube_transition_metrics > container > initialized_to_running_seconds`

**Title:** Pod Initialized to Running

|              |          |
| ------------ | -------- |
| **Type**     | `number` |
| **Required** | No       |

**Description:** The time in seconds from the Pod becoming initialized (all init containers exited 0) to this container running. Only set for non-init containers.

#### <a name="kube_transition_metrics_container_running_timestamp"></a>1.7.5. Property `Metric Record > kube_transition_metrics > container > running_timestamp`

**Title:** Running Timestamp

|              |             |
| ------------ | ----------- |
| **Type**     | `string`    |
| **Required** | No          |
| **Format**   | `date-time` |

**Description:** The timestamp for when the container first entered Running state (first fork(2)/execve(2) in container environment). In the event of a pod restart, this timestamp is NOT updated.

#### <a name="kube_transition_metrics_container_started_timestamp"></a>1.7.6. Property `Metric Record > kube_transition_metrics > container > started_timestamp`

**Title:** Started Timestamp

|              |             |
| ------------ | ----------- |
| **Type**     | `string`    |
| **Required** | No          |
| **Format**   | `date-time` |

**Description:** The timestamp for when the container first started state (startupProbe success). In the event of a pod restart, this timestamp is NOT updated. Only set for non-init containers.

#### <a name="kube_transition_metrics_container_running_to_started_seconds"></a>1.7.7. Property `Metric Record > kube_transition_metrics > container > running_to_started_seconds`

**Title:** Running to Started

|              |          |
| ------------ | -------- |
| **Type**     | `number` |
| **Required** | No       |

**Description:** The time in seconds from the container becoming running to this container started. Only set for non-init containers.

#### <a name="kube_transition_metrics_container_ready_timestamp"></a>1.7.8. Property `Metric Record > kube_transition_metrics > container > ready_timestamp`

**Title:** Started Timestamp

|              |             |
| ------------ | ----------- |
| **Type**     | `string`    |
| **Required** | No          |
| **Format**   | `date-time` |

**Description:** The timestamp for when the container first ready state (readinessProbe success). In the event of a pod restart, this timestamp is NOT updated.

#### <a name="kube_transition_metrics_container_running_to_ready_seconds"></a>1.7.9. Property `Metric Record > kube_transition_metrics > container > running_to_ready_seconds`

**Title:** Running to Ready

|              |          |
| ------------ | -------- |
| **Type**     | `number` |
| **Required** | No       |

**Description:** The time in seconds from the container becoming running to this container ready. In init containers, this is the time the container exited with a successful status.

#### <a name="kube_transition_metrics_container_started_to_ready_seconds"></a>1.7.10. Property `Metric Record > kube_transition_metrics > container > started_to_ready_seconds`

**Title:** Running to Ready

|              |          |
| ------------ | -------- |
| **Type**     | `number` |
| **Required** | No       |

**Description:** The time in seconds from the container becoming started to this container ready. Only set for non-init containers.

### <a name="kube_transition_metrics_image_pull"></a>1.8. Property `Metric Record > kube_transition_metrics > image_pull`

**Title:** Image Pull Metrics

|                           |             |
| ------------------------- | ----------- |
| **Type**                  | `object`    |
| **Required**              | No          |
| **Additional properties** | Not allowed |

**Description:** Included if kube_transition_metric_type is equal to "image_pull". Note that these metrics are only emitted in the event that an image pull occurs, if imagePullPolicy is set to IfNotPresent this will only occur if the image is not already present on the node.

| Property                                                                        | Type    | Title/Description  |
| ------------------------------------------------------------------------------- | ------- | ------------------ |
| + [container_name](#kube_transition_metrics_image_pull_container_name )         | string  | Container name     |
| - [already_present](#kube_transition_metrics_image_pull_already_present )       | boolean | Already Present    |
| + [started_timestamp](#kube_transition_metrics_image_pull_started_timestamp )   | string  | Started Timestamp  |
| - [finished_timestamp](#kube_transition_metrics_image_pull_finished_timestamp ) | string  | Finished Timestamp |
| - [duration_seconds](#kube_transition_metrics_image_pull_duration_seconds )     | number  | Duration           |

#### <a name="kube_transition_metrics_image_pull_container_name"></a>1.8.1. Property `Metric Record > kube_transition_metrics > image_pull > container_name`

**Title:** Container name

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | Yes      |

**Description:** The name of the container which initiated the image pull

#### <a name="kube_transition_metrics_image_pull_already_present"></a>1.8.2. Property `Metric Record > kube_transition_metrics > image_pull > already_present`

**Title:** Already Present

|              |           |
| ------------ | --------- |
| **Type**     | `boolean` |
| **Required** | No        |

**Description:** true if the image was already present on the machine, otherwise false.

#### <a name="kube_transition_metrics_image_pull_started_timestamp"></a>1.8.3. Property `Metric Record > kube_transition_metrics > image_pull > started_timestamp`

**Title:** Started Timestamp

|              |             |
| ------------ | ----------- |
| **Type**     | `string`    |
| **Required** | Yes         |
| **Format**   | `date-time` |

**Description:** The timestamp for when the image pull was first initiated. This is obtained from the Event emitted by the Kubelet and may not be 100% accurate. In the event of ImagePullFailed this time is not reset for subsequent attempts.

#### <a name="kube_transition_metrics_image_pull_finished_timestamp"></a>1.8.4. Property `Metric Record > kube_transition_metrics > image_pull > finished_timestamp`

**Title:** Finished Timestamp

|              |             |
| ------------ | ----------- |
| **Type**     | `string`    |
| **Required** | No          |
| **Format**   | `date-time` |

**Description:** The timestamp for when the image pull was finished. This is obtained from the Event emitted by the Kubelet and may not be 100% accurate.

#### <a name="kube_transition_metrics_image_pull_duration_seconds"></a>1.8.5. Property `Metric Record > kube_transition_metrics > image_pull > duration_seconds`

**Title:** Duration

|              |          |
| ------------ | -------- |
| **Type**     | `number` |
| **Required** | No       |

**Description:** The duration in seconds to complete the image pull successfully. This is based purely off the started_timestamp and finished_timestamp, which themselves are based on Event timestamps which are rounded to seconds. The duration here may not match perfectly the duration seen in the kubelet image pull message, due to slight latency in reporting of image pull Events and truncation of timestamps to seconds.

## <a name="time"></a>2. Property `Metric Record > time`

**Title:** Metric Timestamp

|              |             |
| ------------ | ----------- |
| **Type**     | `string`    |
| **Required** | Yes         |
| **Format**   | `date-time` |

**Description:** The time at which this metric was emitted.

## <a name="message"></a>3. Property `Metric Record > message`

**Title:** Message

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |

**Description:** An additional message emitted along with metrics.

----------------------------------------------------------------------------------------------------------------------------
Generated using [json-schema-for-humans](https://github.com/coveooss/json-schema-for-humans)
