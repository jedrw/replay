package command_test

import (
	"strings"
	"testing"

	"github.com/jedrw/replay/internal/command"
)

func TestParseHistory(t *testing.T) {
	expected := command.ShellHistory{
		command.Command{
			Index:   0,
			Command: "echo 1",
		},
		command.Command{
			Index:   1,
			Command: "echo 2",
		},
		command.Command{
			Index:   2,
			Command: "echo 3",
		},
		command.Command{
			Index:   3,
			Command: "echo 4",
		},
		command.Command{
			Index:   4,
			Command: "echo 5",
		},
	}

	historyString := strings.NewReader("echo 1\necho 2\necho 3\necho 4\necho 5")
	commandHistory, err := command.ParseShellHistory(historyString)
	if err != nil {
		t.Error(err)
	}

	for i, command := range commandHistory {
		if command.Index != expected[i].Index {
			t.Errorf("\nExpected:\n%+v\nGot:\n%+v\n", expected[i].Index, command.Index)
		}
		if command.Command != expected[i].Command {
			t.Errorf("\nExpected:\n%+v\nGot:\n%+v\n", expected[i].Command, command.Command)
		}
	}
}
