package main

import (
	"fmt"
	"os"

	"github.com/jedrw/replay/internal/tui"
)

func main() {
	replayTui := tui.NewReplayTui()

	err := replayTui.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
