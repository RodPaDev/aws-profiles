package state

import (
	"github.com/rodpadev/aws-profiles/lib"
)

type ProfilePosition struct {
	Index int
	Key   string
}

type State struct {
	Selection            ProfilePosition
	ProfileMap           lib.ProfileMap
	ProfileMapSize       int
	ModifiedProfiledKeys []string
}

func (s *State) GetCurrentProfile() (lib.Profile, bool) {
	currentProfile, ok := s.ProfileMap[s.Selection.Key]
	return currentProfile, ok
}
