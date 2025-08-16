package statistics

import (
	"bytes"
	"testing"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/options"
	"github.com/BackMarket-oss/kube-transition-metrics/internal/statistics/state"
	"github.com/BackMarket-oss/kube-transition-metrics/internal/testhelpers"
	"github.com/Izzette/go-safeconcurrency/eventloop"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"
)

func TestNewStatisticEventLoop(t *testing.T) {
	testhelpers.ConfigureLogging(t)

	ctx := t.Context()

	options := &options.Options{
		StatisticEventQueueLength: 1,
		LogLevel:                  zerolog.FatalLevel,
	}

	statisticEventLoop := NewStatisticEventLoop(options)
	defer statisticEventLoop.Close()
	statisticEventLoop.Start()

	gen, err := statisticEventLoop.PodResync(ctx, []types.UID{"test-uid"})
	require.NoError(t, err, "Expected no error when sending resync event")
	state, err := eventloop.WaitForGeneration(ctx, statisticEventLoop, gen)
	require.NoError(t, err, "Expected no error when waiting for resync event to be processed")
	assert.NotNil(t, state, "Expected state to be not nil")

	assert.True(t, state.State().IsBlacklisted("test-uid"), "Expected test-uid to be blacklisted")
	testUIDStatistic, ok := state.State().Get("test-uid")
	assert.Nil(t, testUIDStatistic, "Expected test-uid statistic to be nil")
	assert.False(t, ok, "Expected test-uid statistic to not exist")
	assert.Zero(t, state.State().Len(), "Expected number of tracked statistics to be 0")
}

func TestPodResync(t *testing.T) {
	testhelpers.ConfigureLogging(t)

	buf := &bytes.Buffer{}

	podStatistics := state.NewPodStatistics([]types.UID{})
	assert.False(t, podStatistics.IsBlacklisted("test-uid"), "Expected test-uid to not be blacklisted")
	assert.Zero(t, podStatistics.Len(), "Expected number of tracked statistics to be 0")

	ev := &resyncEvent{
		blacklistUIDs: []types.UID{"test-uid"},
		output:        buf,
	}
	nextPodStatistics := ev.Dispatch(0, podStatistics)
	assert.NotNil(t, nextPodStatistics, "Expected nextPodStatistics to be not nil")
	assert.True(t, nextPodStatistics.IsBlacklisted("test-uid"), "Expected test-uid to be blacklisted")
	assert.Zero(t, nextPodStatistics.Len(), "Expected number of tracked statistics to be 0")
	assert.Empty(t, buf.String(), "Expected no output")
}
