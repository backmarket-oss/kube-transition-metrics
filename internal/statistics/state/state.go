package state

import (
	"github.com/Izzette/go-safeconcurrency/eventloop/snapshot"
	"github.com/benbjohnson/immutable"
	"k8s.io/apimachinery/pkg/types"
)

// State holds the statistics for the pods and image pulls.
// State is immutable, all the methods return a new instance of the struct.
// Do not lose track of the returned instance, it should be assigned to the containing structure.
type State struct {
	podStatistics       *PodStatistics
	imagePullStatistics *immutable.Map[types.UID, *PodImagePullStatistic]
}

// NewState creates a new State instance with empty statistics.
func NewState() *State {
	return &State{
		podStatistics:       NewPodStatistics([]types.UID{}),
		imagePullStatistics: &immutable.Map[types.UID, *PodImagePullStatistic]{},
	}
}

// Copy implements [github.com/Izzette/go-safeconcurrency/types.Copyable.Copy].
func (s *State) Copy() *State {
	return snapshot.CopyPtr(s)
}

// GetPodStatistics returns the pod statistics.
func (s *State) GetPodStatistics() *PodStatistics {
	return s.podStatistics
}

// SetPodStatistics sets the pod statistics and returns a new State instance.
// SetPodStatistics returns a new instance of the State with the updated fields.
func (s *State) SetPodStatistics(podStatistics *PodStatistics) *State {
	s = s.Copy()
	s.podStatistics = podStatistics

	return s
}

// LenImagePullStatistics returns the number of image pull statistics.
func (s *State) LenImagePullStatistics() int {
	return s.imagePullStatistics.Len()
}

// GetImagePullStatistic returns the image pull statistic for the pod UID, if it exists.
func (s *State) GetImagePullStatistic(uid types.UID) (*PodImagePullStatistic, bool) {
	return s.imagePullStatistics.Get(uid)
}

// SetImagePullStatistic sets the image pull statistic for the given pod UID.
// SetImagePullStatistic returns a new instance of the State with the updated fields.
func (s *State) SetImagePullStatistic(uid types.UID, imagePullStatistic *PodImagePullStatistic) *State {
	s = s.Copy()
	s.imagePullStatistics = s.imagePullStatistics.Set(uid, imagePullStatistic)

	return s
}

// DeleteImagePullStatistic deletes the image pull statistic for the given pod UID, if it exists.
// DeleteImagePullStatistic returns a new instance of the State with the updated fields.
func (s *State) DeleteImagePullStatistic(uid types.UID) *State {
	s = s.Copy()
	s.imagePullStatistics = s.imagePullStatistics.Delete(uid)

	return s
}
