package state

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestNewPodStatistics(t *testing.T) {
	state := NewPodStatistics([]types.UID{})
	assert.Equal(t, 0, state.Len(), "Expected length to be 0")
	pod, ok := state.Get("test-uid")
	assert.False(t, ok, "Expected to not find pod with UID test-uid")
	assert.Nil(t, pod, "Expected pod to be nil")
}

func TestPodStatisticsSetAndGet(t *testing.T) {
	state := NewPodStatistics([]types.UID{})
	uid := types.UID("test-uid")
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "test-namespace",
			UID:       uid,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "test-container"},
			},
		},
		Status: corev1.PodStatus{
			Conditions: []corev1.PodCondition{
				{
					Type:   corev1.PodReady,
					Status: corev1.ConditionTrue,
				},
				{
					Type:   corev1.PodInitialized,
					Status: corev1.ConditionTrue,
				},
			},
		},
	}

	podStat := NewPodStatistic(time.Now(), pod)
	state = state.Set(uid, podStat)
	assert.Equal(t, 1, state.Len(), "Expected length to be 1")

	retrievedPodStat, ok := state.Get(uid)
	assert.True(t, ok, "Expected to find pod with UID test-uid")
	assert.Equal(t, podStat, retrievedPodStat, "Expected retrieved pod statistic to match set pod statistic")
}

func TestPodStatisticsDelete(t *testing.T) {
	state := NewPodStatistics([]types.UID{})
	uid := types.UID("test-uid")
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "test-namespace",
			UID:       uid,
		},
	}

	podStat := NewPodStatistic(time.Now(), pod)
	state = state.Set(uid, podStat)
	assert.Equal(t, 1, state.Len(), "Expected length to be 1")

	state = state.Delete(uid)
	assert.Equal(t, 0, state.Len(), "Expected length to be 0")
}

func TestPodStatisticsIsBlacklisted(t *testing.T) {
	uid := types.UID("test-uid")
	state := NewPodStatistics([]types.UID{uid})

	assert.True(t, state.IsBlacklisted(uid), "Expected UID to be blacklisted")
	assert.False(t, state.IsBlacklisted(types.UID("other-uid")), "Expected other UID to not be blacklisted")
}
