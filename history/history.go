package history

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
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

func historyFilePath() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Could not find history file path: %s", err)
	}

	return path.Join(homedir, fmt.Sprintf(".%s_history", GetShell()))
}

func getHistoryFile(path string) (*os.File, error) {
	return os.Open(path)
}

func parseHistory(historyFile io.Reader) (CommandHistory, error) {
	var commandHistory CommandHistory
	reader := bufio.NewReader(historyFile)
	i := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			return commandHistory, err
		}

		commandHistory = append(
			commandHistory,
			Command{
				Index:   i,
				Command: strings.TrimRight(line, "\n"),
			},
		)

		i++
	}

	return commandHistory, nil
}

func GetHistory() (CommandHistory, error) {
	historyFile, err := getHistoryFile(historyFilePath())
	if err != nil {
		return CommandHistory{}, err
	}
	defer historyFile.Close()

	return parseHistory(historyFile)
}
