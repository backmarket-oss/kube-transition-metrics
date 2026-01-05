package state

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func p[T any](v T) *T {
	return &v
}

func testLogger(t *testing.T) (*zerolog.Logger, *bytes.Buffer) {
	t.Helper()

	buf := &bytes.Buffer{}
	logger := zerolog.New(buf)

	return &logger, buf
}

func parseEventBuffer(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()

	var event map[string]any

	err := json.Unmarshal(buf.Bytes(), &event)
	require.NoError(t, err, "Failed to unmarshal JSON")

	return event
}

func TestCommonPodLabelsAll(t *testing.T) {
	// Define a test pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "test-namespace",
			Labels: map[string]string{
				"app.kubernetes.io/component":  "test-component",
				"app.kubernetes.io/instance":   "test-instance",
				"app.kubernetes.io/managed-by": "test-managed-by",
				"app.kubernetes.io/name":       "test-name",
				"app.kubernetes.io/part-of":    "test-part-of",
				"app.kubernetes.io/version":    "test-version",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "apps/v1",
					Kind:       "ReplicaSet",
					Name:       "test-replicaset",
					Controller: p(true),
				},
			},
		},
		Spec: corev1.PodSpec{
			NodeName:          "test-node",
			PriorityClassName: "test-priority-class",
			RuntimeClassName:  p("test-runtime-class"),
		},
		Status: corev1.PodStatus{
			QOSClass: corev1.PodQOSGuaranteed,
		},
	}

	// Get the logger function
	commonLabels := commonPodLabels(pod)

	// Create a zerolog event
	logger, buf := testLogger(t)
	logger.Info().Func(commonLabels).Msg("Test message")

	// Assertions
	event := parseEventBuffer(t, buf)
	assert.Equal(t, map[string]any{
		"kube_namespace":      pod.Namespace,
		"pod_name":            pod.Name,
		"kube_node":           pod.Spec.NodeName,
		"kube_qos":            string(pod.Status.QOSClass),
		"kube_priority_class": pod.Spec.PriorityClassName,
		"kube_runtime_class":  *pod.Spec.RuntimeClassName,
		"kube_app_component":  "test-component",
		"kube_app_instance":   "test-instance",
		"kube_app_managed_by": "test-managed-by",
		"kube_app_name":       "test-name",
		"kube_app_part_of":    "test-part-of",
		"kube_app_version":    "test-version",
		"kube_ownerref_kind":  "replicaset",
		"kube_ownerref_name":  "test-replicaset",
		"kube_replica_set":    "test-replicaset",
		"message":             "Test message",
		"level":               "info",
	}, event)
}

func TestCommonPodLabelsNode(t *testing.T) {
	t.Parallel()

	for _, test := range []commonNodeNameTest{
		{
			name:      "With",
			nodeName:  "test-node",
			expectNil: false,
		},
		{
			name:      "Without",
			nodeName:  "",
			expectNil: true,
		},
	} {
		t.Run(test.name, test.Test)
	}
}

type commonNodeNameTest struct {
	name      string
	nodeName  string
	expectNil bool
}

func (st *commonNodeNameTest) Test(t *testing.T) {
	t.Parallel()

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "test-namespace",
		},
		Spec: corev1.PodSpec{
			NodeName: st.nodeName,
		},
	}

	commonLabels := commonPodLabels(pod)

	logger, buf := testLogger(t)
	logger.Info().Func(commonLabels).Msg("Test message")

	// Assertions
	event := parseEventBuffer(t, buf)
	if st.expectNil {
		assert.NotContains(t, event, "kube_node")
	} else {
		assert.Contains(t, event, "kube_node")
		assert.Equal(t, event["kube_node"], st.nodeName)
	}
}

func TestCommonPodLabelsQOS(t *testing.T) {
	t.Parallel()

	for _, test := range []commonQOSTest{
		{
			name:      "With",
			qosClass:  corev1.PodQOSBurstable,
			expectNil: false,
		},
		{
			name:      "Without",
			qosClass:  "",
			expectNil: true,
		},
	} {
		t.Run(test.name, test.Test)
	}
}

type commonQOSTest struct {
	name      string
	qosClass  corev1.PodQOSClass
	expectNil bool
}

func (st *commonQOSTest) Test(t *testing.T) {
	t.Parallel()

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "test-namespace",
		},
		Status: corev1.PodStatus{
			QOSClass: st.qosClass,
		},
	}

	commonLabels := commonPodLabels(pod)

	logger, buf := testLogger(t)
	logger.Info().Func(commonLabels).Msg("Test message")

	event := parseEventBuffer(t, buf)

	if st.expectNil {
		assert.NotContains(t, event, "kube_qos")
	} else {
		assert.Contains(t, event, "kube_qos")
		assert.Equal(t, event["kube_qos"], string(st.qosClass))
	}
}

func TestCommonPodLabelsPriorityClass(t *testing.T) {
	t.Parallel()

	for _, test := range []commonPriorityClassTest{
		{
			name:      "With",
			className: "test-priority-class",
			expectNil: false,
		},
		{
			name:      "Without",
			className: "",
			expectNil: true,
		},
	} {
		t.Run(test.name, test.Test)
	}
}

type commonPriorityClassTest struct {
	name      string
	className string
	expectNil bool
}

func (st *commonPriorityClassTest) Test(t *testing.T) {
	t.Parallel()

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "test-namespace",
		},
		Spec: corev1.PodSpec{
			PriorityClassName: st.className,
		},
	}

	commonLabels := commonPodLabels(pod)

	logger, buf := testLogger(t)
	logger.Info().Func(commonLabels).Msg("Test message")

	event := parseEventBuffer(t, buf)

	if st.expectNil {
		assert.NotContains(t, event, "kube_priority_class")
	} else {
		assert.Contains(t, event, "kube_priority_class")
		assert.Equal(t, event["kube_priority_class"], st.className)
	}
}

func TestCommonPodLabelsRuntimeClass(t *testing.T) {
	t.Parallel()

	for _, test := range []commonRuntimeClassTest{
		{
			name:      "With",
			className: p("test-runtime-class"),
			expectNil: false,
		},
		{
			name:      "WithEmpty",
			className: p(""), // not sure if an empty string is valid, but the event should probably emit what is there.
			expectNil: false,
		},
		{
			name:      "Without",
			className: nil,
			expectNil: true,
		},
	} {
		t.Run(test.name, test.Test)
	}
}

type commonRuntimeClassTest struct {
	name      string
	className *string
	expectNil bool
}

func (st *commonRuntimeClassTest) Test(t *testing.T) {
	t.Parallel()

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "test-namespace",
		},
		Spec: corev1.PodSpec{
			RuntimeClassName: st.className,
		},
	}

	commonLabels := commonPodLabels(pod)

	logger, buf := testLogger(t)
	logger.Info().Func(commonLabels).Msg("Test message")

	event := parseEventBuffer(t, buf)

	if st.expectNil {
		assert.NotContains(t, event, "kube_runtime_class")
	} else {
		assert.Contains(t, event, "kube_runtime_class")
		assert.Equal(t, event["kube_runtime_class"], *st.className)
	}
}

func TestOwnerRefLabels(t *testing.T) {
	t.Parallel()

	for _, test := range []ownerRefLabelsTest{
		{
			name: "WithReplicaSet",
			ownerRefs: []metav1.OwnerReference{
				{
					APIVersion: "apps/v1",
					Kind:       "ReplicaSet",
					Name:       "test-replicaset",
					Controller: p(true),
				},
			},
			expectNil:        false,
			expectKind:       "replicaset",
			expectName:       "test-replicaset",
			expectOtherField: "kube_replica_set",
			expectOtherValue: "test-replicaset",
		},
		{
			name: "WithDaemonSet",
			ownerRefs: []metav1.OwnerReference{
				{
					APIVersion: "apps/v1",
					Kind:       "DaemonSet",
					Name:       "test-daemonset",
					Controller: p(true),
				},
			},
			expectNil:        false,
			expectKind:       "daemonset",
			expectName:       "test-daemonset",
			expectOtherField: "kube_daemon_set",
			expectOtherValue: "test-daemonset",
		},
		{
			name: "WithStatefulSet",
			ownerRefs: []metav1.OwnerReference{
				{
					APIVersion: "apps/v1",
					Kind:       "StatefulSet",
					Name:       "test-statefulset",
					Controller: p(true),
				},
			},
			expectNil:  false,
			expectKind: "statefulset",
			expectName: "test-statefulset",
		},
		{
			name: "WithJob",
			ownerRefs: []metav1.OwnerReference{
				{
					APIVersion: "batch/v1",
					Kind:       "Job",
					Name:       "test-job",
					Controller: p(true),
				},
			},
			expectNil:  false,
			expectKind: "job",
			expectName: "test-job",
		},
		{
			name: "WithNonController",
			ownerRefs: []metav1.OwnerReference{
				{
					APIVersion: "custom.io/v1",
					Kind:       "CustomResource",
					Name:       "test-custom-resource",
					Controller: p(false),
				},
			},
			expectNil: true,
		},
		{
			name: "WithMultiple",
			ownerRefs: []metav1.OwnerReference{
				{
					APIVersion: "custom.io/v1",
					Kind:       "CustomResource",
					Name:       "test-custom-resource",
					Controller: p(false),
				},
				{
					APIVersion: "apps/v1",
					Kind:       "ReplicaSet",
					Name:       "test-replicaset",
					Controller: p(true),
				},
			},
			expectNil:        false,
			expectKind:       "replicaset",
			expectName:       "test-replicaset",
			expectOtherField: "kube_replica_set",
			expectOtherValue: "test-replicaset",
		},
		{
			name: "WithCustom",
			ownerRefs: []metav1.OwnerReference{
				{
					APIVersion: "custom.io/v1",
					Kind:       "CustomResource",
					Name:       "test-custom-resource",
					Controller: p(true),
				},
			},
			expectNil:  false,
			expectKind: "customresource",
			expectName: "test-custom-resource",
		},
		{
			name:      "Without",
			ownerRefs: []metav1.OwnerReference{},
			expectNil: true,
		},
	} {
		t.Run(test.name, test.Test)
	}
}

type ownerRefLabelsTest struct {
	name       string
	ownerRefs  []metav1.OwnerReference
	expectNil  bool
	expectKind string
	expectName string

	expectOtherField string
	expectOtherValue string
}

func (st *ownerRefLabelsTest) Test(t *testing.T) {
	t.Parallel()

	ownerRefLabels := ownerRefLabels(st.ownerRefs)

	logger, buf := testLogger(t)
	logger.Info().Func(ownerRefLabels).Msg("Test message")

	event := parseEventBuffer(t, buf)

	if st.expectNil {
		assert.NotContains(t, event, "kube_ownerref_kind")
		assert.NotContains(t, event, "kube_ownerref_name")
	} else {
		assert.Contains(t, event, "kube_ownerref_kind")
		assert.Contains(t, event, "kube_ownerref_name")
		assert.Equal(t, event["kube_ownerref_kind"], st.expectKind)
		assert.Equal(t, event["kube_ownerref_name"], st.expectName)
	}

	if st.expectOtherField != "" {
		assert.Contains(t, event, st.expectOtherField)
		assert.Equal(t, event[st.expectOtherField], st.expectOtherValue)
	}
}

func TestAppLables(t *testing.T) {
	t.Parallel()

	for _, test := range []appLabelsTest{
		{
			name: "WithAll",
			labels: map[string]string{
				"app.kubernetes.io/component":  "test-component",
				"app.kubernetes.io/instance":   "test-instance",
				"app.kubernetes.io/managed-by": "test-managed-by",
				"app.kubernetes.io/name":       "test-name",
				"app.kubernetes.io/part-of":    "test-part-of",
				"app.kubernetes.io/version":    "test-version",
			},
			expectedLabels: map[string]string{
				"kube_app_component":  "test-component",
				"kube_app_instance":   "test-instance",
				"kube_app_managed_by": "test-managed-by",
				"kube_app_name":       "test-name",
				"kube_app_part_of":    "test-part-of",
				"kube_app_version":    "test-version",
			},
		},
		{
			name: "WithSome",
			labels: map[string]string{
				"app.kubernetes.io/component":  "test-component",
				"app.kubernetes.io/instance":   "test-instance",
				"app.kubernetes.io/managed-by": "test-managed-by",
			},
			expectedLabels: map[string]string{
				"kube_app_component":  "test-component",
				"kube_app_instance":   "test-instance",
				"kube_app_managed_by": "test-managed-by",
			},
		},
		{
			name:           "Without",
			labels:         map[string]string{},
			expectedLabels: map[string]string{},
		},
		{
			name: "WithUnknown",
			labels: map[string]string{
				"app.kubernetes.io/invalid": "test-invalid",
			},
			expectedLabels: map[string]string{},
		},
	} {
		t.Run(test.name, test.Test)
	}
}

type appLabelsTest struct {
	name           string
	labels         map[string]string
	expectedLabels map[string]string
}

func (st *appLabelsTest) Test(t *testing.T) {
	t.Parallel()

	appLabels := appLabels(st.labels)

	logger, buf := testLogger(t)
	logger.Info().Func(appLabels).Msg("Test message")

	event := parseEventBuffer(t, buf)

	for key := range appLabelFields {
		if expectedValue, ok := st.expectedLabels[key]; ok {
			assert.Contains(t, event, key)
			assert.Equal(t, expectedValue, event[key])
		} else {
			assert.NotContains(t, event, key)
		}
	}
}

func TestCommonContainerLabelsAll(t *testing.T) {
	t.Parallel()

	// Define a test container
	container := &corev1.Container{
		Name:  "test-container",
		Image: "docker.io/library/test-image:latest",
	}

	errorLogger, errorBuf := testLogger(t)
	commonLabels := commonContainerLabels(errorLogger, container)

	// Create a zerolog event
	logger, buf := testLogger(t)
	logger.Info().Func(commonLabels).Msg("Test message")

	// Assertions
	assert.Empty(t, errorBuf.String())
	event := parseEventBuffer(t, buf)
	assert.Equal(t, map[string]any{
		"container_name": "test-container",
		"short_image":    "test-image",
		"image_name":     "docker.io/library/test-image",
		"image_tag":      "latest",
		"level":          "info",
		"message":        "Test message",
	}, event)
}

func TestImageLabels(t *testing.T) {
	t.Parallel()

	for _, test := range []imageLabelsTest{
		{
			name:             "WithFull",
			image:            "docker.io/library/test-image:latest",
			expectImageName:  "docker.io/library/test-image",
			expectImageTag:   "latest",
			expectShortImage: "test-image",
		},
		{
			name:             "WithTag",
			image:            "docker.io/library/test-image:xyz",
			expectImageName:  "docker.io/library/test-image",
			expectImageTag:   "xyz",
			expectShortImage: "test-image",
		},
		{
			name: "WithDigest",
			image: "docker.io/library/test-image@sha256:" +
				"01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b",
			expectImageName:  "docker.io/library/test-image",
			expectImageTag:   "sha256:01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b",
			expectShortImage: "test-image",
		},
		{
			name: "WithTagAndDigest",
			image: "docker.io/library/test-image:latest@sha256:" +
				"01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b",
			expectImageName:  "docker.io/library/test-image",
			expectImageTag:   "latest",
			expectShortImage: "test-image",
		},
		{
			name:             "WithImplicitRegistry",
			image:            "test-image:latest",
			expectImageName:  "docker.io/library/test-image",
			expectImageTag:   "latest",
			expectShortImage: "test-image",
		},
		{
			name:             "WithImplicitTag",
			image:            "docker.io/library/test-image",
			expectImageName:  "docker.io/library/test-image",
			expectImageTag:   "latest",
			expectShortImage: "test-image",
		},
		{
			name:        "WithInvalidImage",
			image:       "invalid-image@sha256:0", // 0 is not a valid digest
			expectError: "failed to parse image name",
		},
		{
			name: "WithEncoded",
			// URI encoding is not allowed in image names (who knows why)
			image:       "fancy.com:6000/my%2dstuff/test-image:latest",
			expectError: "failed to parse image name",
		},
	} {
		t.Run(test.name, test.Test)
	}
}

type imageLabelsTest struct {
	name             string
	image            string
	expectError      string
	expectImageName  string
	expectImageTag   string
	expectShortImage string
}

func (st *imageLabelsTest) Test(t *testing.T) {
	t.Parallel()

	errorLogger, errorBuf := testLogger(t)
	imageLabels := imageLabels(errorLogger, st.image)

	logger, buf := testLogger(t)
	logger.Info().Func(imageLabels).Msg("Test message")

	if st.expectError != "" {
		errors := parseEventBuffer(t, errorBuf)
		assert.Equal(t, errors["message"], st.expectError)

		// if there is an error, there should be no event
		assert.NotContains(t, errors, "image_name")
		assert.NotContains(t, errors, "image_tag")
		assert.NotContains(t, errors, "short_image")

		return
	} else {
		assert.Empty(t, errorBuf.String())
	}

	event := parseEventBuffer(t, buf)

	assert.Contains(t, event, "image_name")
	assert.Equal(t, event["image_name"], st.expectImageName)

	assert.Contains(t, event, "image_tag")
	assert.Equal(t, event["image_tag"], st.expectImageTag)

	assert.Contains(t, event, "short_image")
	assert.Equal(t, event["short_image"], st.expectShortImage)
}
