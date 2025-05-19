package state

import (
	"github.com/rodpadev/aws-profiles/lib"
)

type ProfilePosition struct {
	Index int
	Key   string
}

type State struct {
	Cursor     ProfilePosition
	Selection  ProfilePosition
	ProfileMap lib.ProfileMap
}

func (s *State) GetCurrentProfile() (lib.Profile, bool) {
	currentProfile, ok := s.ProfileMap[s.Selection.Key]
	return currentProfile, ok
}
