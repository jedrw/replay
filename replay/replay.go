package replay

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jedrw/replay/command"
)

func newReplayCommand(command command.Command) *exec.Cmd {
	replayCmd := exec.Command("bash", "-c")
	replayCmd.Args = append(replayCmd.Args, command.Command)
	replayCmd.Stdin = os.Stdin
	replayCmd.Stdout = os.Stdout
	replayCmd.Stderr = os.Stderr

	return replayCmd
}

func Replay(commands []command.Command) {
	for _, command := range commands {
		replayCmd := newReplayCommand(command)
		err := replayCmd.Start()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := replayCmd.Wait(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok && exiterr.ExitCode() != 0 {
				os.Exit(0)
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}
