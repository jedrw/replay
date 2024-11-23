package history_test

import (
	"reflect"
	"testing"

	"github.com/jedrw/replay/history"
)

var testReplayHistory = history.ReplayHistory{
	history.Replay{
		"ls -la",
		"echo pickle",
	},
	history.Replay{
		"echo sausage",
		"ls -la",
	},
	history.Replay{
		"ls -l",
		"echo mustard",
	},
	history.Replay{
		"cat dog",
	},
}

func TestUpdateReplayHistoryPrependsNewReplay(t *testing.T) {
	newReplay := history.Replay{
		"ls -la",
		"echo cucumber",
	}

	newHistory := history.UpdateReplayHistory(newReplay, testReplayHistory)

	expectedHistory := append(
		history.ReplayHistory{newReplay},
		testReplayHistory...,
	)

	if !reflect.DeepEqual(newHistory, expectedHistory) {
		t.Errorf("\nexpected: %+v\n     got: %+v", expectedHistory, newHistory)
	}
}

func TestUpdateReplayHistoryRemovesDuplicateReplay(t *testing.T) {
	newReplay := history.Replay{
		"cat dog",
	}

	newHistory := history.UpdateReplayHistory(newReplay, testReplayHistory)

	expectedHistory := history.ReplayHistory{
		history.Replay{
			"cat dog",
		},
		history.Replay{
			"ls -la",
			"echo pickle",
		},
		history.Replay{
			"echo sausage",
			"ls -la",
		},
		history.Replay{
			"ls -l",
			"echo mustard",
		},
	}

	if !reflect.DeepEqual(newHistory, expectedHistory) {
		t.Errorf("\nexpected: %+v\n     got: %+v", expectedHistory, newHistory)
	}
}
