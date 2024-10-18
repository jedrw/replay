package history_test

import (
	"strings"
	"testing"

	"github.com/lupinelab/replay/history"
)

func TestParseHistory(t *testing.T) {
	expected := history.CommandHistory{
		history.Command{
			Index:   0,
			Command: "echo 1",
		},
		history.Command{
			Index:   1,
			Command: "echo 2",
		},
		history.Command{
			Index:   2,
			Command: "echo 3",
		},
		history.Command{
			Index:   3,
			Command: "echo 4",
		},
		history.Command{
			Index:   4,
			Command: "echo 5",
		},
	}

	historyString := strings.NewReader("echo 1\necho 2\necho 3\necho 4\necho 5")
	commandHistory, err := history.ParseHistory(historyString)
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
