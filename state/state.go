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
