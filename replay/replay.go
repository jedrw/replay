package replay

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/lupinelab/replay/history"
)

func newReplayCommand(command history.Command) *exec.Cmd {
	replayCmd := exec.Command("bash", "-c")
	replayCmd.Args = append(replayCmd.Args, command.Command)
	replayCmd.Stdin = os.Stdin
	replayCmd.Stdout = os.Stdout
	replayCmd.Stderr = os.Stderr

	return replayCmd
}

func Replay(commands []history.Command) {
	for _, command := range commands {
		replayCmd := newReplayCommand(command)
		err := replayCmd.Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
