package command

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

type ShellHistory []Command

func getShell() string {
	shellBin := os.Getenv("SHELL")
	return path.Base(shellBin)
}

func ShellHistoryPath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(homedir, fmt.Sprintf(".%s_history", getShell())), nil
}

func parseShellHistory(historyFile io.Reader) (ShellHistory, error) {
	var history ShellHistory
	reader := bufio.NewReader(historyFile)
	i := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			return history, err
		}

		history = append(
			history,
			Command{
				Index:   i,
				Command: strings.TrimRight(line, "\n"),
			},
		)

		i++
	}

	return history, nil
}

func GetShellHistory(historyPath string) (ShellHistory, error) {
	historyFile, err := os.Open(historyPath)
	if err != nil {
		return ShellHistory{}, err
	}
	defer historyFile.Close()

	return parseShellHistory(historyFile)
}
