package state

import (
	"github.com/rodpadev/aws-profiles/lib"
)

type ProfilePosition struct {
	Index int
	Key   string
}

type State struct {
	Selection           ProfilePosition
	ProfileMap          lib.ProfileMap
	EditedProfileMap    lib.ProfileMap
	IsCommandModeActive bool
}

func (s *State) GetCurrentProfile() (lib.Profile, bool) {
	currentProfile, ok := s.ProfileMap[s.Selection.Key]
	return currentProfile, ok
}

func (s *State) IsFieldEdited(profileKey string, field string) bool {
	if s.EditedProfileMap == nil {
		return false
	}
	editedProfile, ok := s.EditedProfileMap[profileKey]
	if !ok {
		return false
	}
	originalProfile, ok := s.ProfileMap[profileKey]
	if !ok {
		return false
	}
	originalField := originalProfile.GetField(field)
	editedField := editedProfile.GetField(field)

	return originalField != editedField

}
