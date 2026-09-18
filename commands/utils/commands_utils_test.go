package utils

import (
	"patrol_build_flutter/commands"
	"testing"
)

func Test_IsSameCommand(t *testing.T) {
	t.Run("correct command returns true", func(t *testing.T) {
		correctFirst := commands.FlutterPubDependencies
		correctSecond := commands.FlutterPubDependencies
		isSameCommand := IsSameCommand(correctFirst, correctSecond)
		if !isSameCommand {
			t.Error("expected true for identical commands, got false")
		}
	})

	t.Run("wrong command returns false", func(t *testing.T) {
		correctCommand := commands.FlutterPubDependencies
		wrongCmd := commands.Command{Name: "echo", Args: []string{"hello"}}
		isSameCommand := IsSameCommand(wrongCmd, correctCommand)
		if isSameCommand {
			t.Error("expected false for different commands, got true")
		}
	})

	t.Run("different args returns false", func(t *testing.T) {
		cmd1 := commands.Command{Name: "flutter", Args: []string{"pub", "dependencies"}}
		cmd2 := commands.Command{Name: "flutter", Args: []string{"pub", "get"}}
		isSameCommand := IsSameCommand(cmd1, cmd2)
		if isSameCommand {
			t.Error("expected false for commands with different args, got true")
		}
	})
	t.Run("different args length", func(t *testing.T) {
		cmd1 := commands.Command{Name: "flutter", Args: []string{"pub", "get"}}
		cmd2 := commands.Command{Name: "flutter", Args: []string{"build", "apk", "--release"}}
		isSameCommand := IsSameCommand(cmd1, cmd2)
		if isSameCommand {
			t.Error("expected false for commands with different args length, got true")
		}
	})
}
