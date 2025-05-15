package statistics

import (
	"testing"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/logging"
	"github.com/BackMarket-oss/kube-transition-metrics/internal/options"
	"github.com/Izzette/go-safeconcurrency/eventloop"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"
)

func TestNewStatisticEventHandler(t *testing.T) {
	ctx := t.Context()

	options := &options.Options{
		StatisticEventQueueLength: 1,
		LogLevel:                  zerolog.FatalLevel,
	}
	logging.Configure()
	logging.SetOptions(options)

	statisticsEventLoop := NewStatisticEventLoop(options)
	defer statisticsEventLoop.Close()
	statisticsEventLoop.Start()

	resyncEvent := &resyncEvent{
		blacklistUIDs: []types.UID{"test-uid"},
	}

	state, err := eventloop.SendAndWait(ctx, statisticsEventLoop, resyncEvent)
	require.NoError(t, err, "Expected no error when sending resync event")
	assert.NotNil(t, state, "Expected state to be not nil")
	assert.True(t, state.State().IsBlacklisted("test-uid"), "Expected test-uid to be blacklisted")
	testUIDStatistic, ok := state.State().Get("test-uid")
	assert.Nil(t, testUIDStatistic, "Expected test-uid statistic to be nil")
	assert.False(t, ok, "Expected test-uid statistic to not exist")
	assert.Zero(t, state.State().Len(), "Expected number of tracked statistics to be 0")
}
