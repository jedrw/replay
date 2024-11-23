package history

import (
	"encoding/json"
	"io"
	"math"
	"os"
	"reflect"

	"github.com/adrg/xdg"
	"github.com/jedrw/replay/command"
)

type Replay []string
type ReplayHistory []Replay

func NewReplayFromCommands(commands []command.Command) Replay {
	replay := Replay{}
	for _, c := range commands {
		replay = append(replay, c.Command)
	}

	return replay
}

func ReplayHistoryPath() (string, error) {
	historyPath, err := xdg.DataFile("replay/history.json")
	if err != nil {
		return "", err
	}

	return historyPath, nil
}

func GetReplayHistory(historyPath string) (ReplayHistory, error) {
	historyFile, err := os.OpenFile(historyPath, os.O_RDONLY|os.O_CREATE, 0640)
	if err != nil {
		return ReplayHistory{}, err
	}
	defer historyFile.Close()

	historyBytes, err := io.ReadAll(historyFile)
	if err != nil {
		return ReplayHistory{}, err
	}

	if len(historyBytes) == 0 {
		return ReplayHistory{}, nil
	}

	var history ReplayHistory
	err = json.Unmarshal(historyBytes, &history)
	if err != nil {
		return ReplayHistory{}, err
	}

	return history, nil
}

func UpdateReplayHistory(replay Replay, history ReplayHistory) ReplayHistory {
	if len(replay) == 0 {
		return history
	} else {
		newHistory := ReplayHistory{replay}
		for _, r := range history {
			if !(reflect.DeepEqual(r, replay) || len(r) == 0) {
				newHistory = append(newHistory, r)
			}
		}

		return newHistory[:int(math.Min(10, float64(len(newHistory))))]
	}
}

func WriteReplayHistory(history ReplayHistory, historyPath string) error {
	historyBytes, err := json.Marshal(history)
	if err != nil {
		return err
	}

	historyFile, err := os.OpenFile(historyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		return err
	}
	defer historyFile.Close()

	_, err = historyFile.Write(historyBytes)
	if err != nil {
		return err
	}

	return nil
}
