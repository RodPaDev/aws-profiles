package state

import (
	"github.com/rodpadev/aws-profiles/lib"
)

type Selection struct {
	Index int
	Key   string
}

type State struct {
	Selection  Selection
	ProfileMap lib.ProfileMap
}
