package utils

import (
	"patrol_build_flutter/commands"
)

// isSameCommand compares two commands.Command structs by Name and Args.
func IsSameCommand(a, b commands.Command) bool {
	if a.Name != b.Name {
		return false
	}
	if len(a.Args) != len(b.Args) {
		return false
	}
	for i := range a.Args {
		if a.Args[i] != b.Args[i] {
			return false
		}
	}
	return true
}
