package history

import (
	"strings"
	"testing"
)

func TestParseHistory(t *testing.T) {
	expected := CommandHistory{
		Command{
			Index:   0,
			Command: "echo 1",
		},
		Command{
			Index:   1,
			Command: "echo 2",
		},
		Command{
			Index:   2,
			Command: "echo 3",
		},
		Command{
			Index:   3,
			Command: "echo 4",
		},
		Command{
			Index:   4,
			Command: "echo 5",
		},
	}

	history := strings.NewReader("echo 1\necho 2\necho 3\necho 4\necho 5")
	commandHistory, err := parseHistory(history)
	if err != nil {
		t.Error(err)
	}

	for i, command := range commandHistory {
		if command.Index != expected[i].Index {
			t.Errorf("\nExpected:\n%+v\nGot:\n%+v", expected[i].Index, command.Index)
		}
		if command.Command != expected[i].Command {
			t.Errorf("\nExpected:\n%+v\nGot:\n%+v", expected[i].Command, command.Command)
		}
	}
}
