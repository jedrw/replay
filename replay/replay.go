package replay

import (
	"os"
	"os/exec"

	"github.com/lupinelab/replay/history"
)

func NewReplay(commands []history.Command) *exec.Cmd {
	replayCmd := exec.Command("bash", "-c")

	var replayString string
	for i, command := range commands {
		replayString += command.Command
		if i < len(commands)-1 {
			replayString += " && "
		}
	}

	replayCmd.Args = append(replayCmd.Args, replayString)
	replayCmd.Stdin = os.Stdin
	replayCmd.Stdout = os.Stdout
	replayCmd.Stderr = os.Stderr

	return replayCmd
}
