package history

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"slices"
)

type Command struct {
	Index   int
	Command string
}

type CommandHistory []Command

func GetShell() string {
	shellBin := os.Getenv("SHELL")
	return path.Base(shellBin)
}

func readHistory() ([]string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	historyFile, err := os.Open(path.Join(homedir, fmt.Sprintf(".%s_history", GetShell())))
	if err != nil {
		return nil, err
	}

	defer historyFile.Close()

	var historyLines []string
	scanner := bufio.NewScanner(historyFile)
	for scanner.Scan() {
		historyLines = append(historyLines, scanner.Text())
	}

	return historyLines, scanner.Err()
}

func parseHistory(historyLines []string) (CommandHistory, error) {
	var commandHistory CommandHistory
	for i, line := range historyLines {
		commandHistory = append(
			CommandHistory{
				Command{
					Index:   i,
					Command: line,
				},
			},
			commandHistory...,
		)
	}

	slices.Reverse(commandHistory)
	return commandHistory, nil
}

func GetHistory() (CommandHistory, error) {
	historyBytes, err := readHistory()
	if err != nil {
		return CommandHistory{}, err
	}

	return parseHistory(historyBytes)
}
