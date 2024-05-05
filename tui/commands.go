package tui

import (
	"cmp"
	"slices"

	"github.com/lupinelab/replay/history"
)

func sortCommands(commands []command) []history.Command {
	var orderedCommands []command
	var remainingCommands []history.Command
	for _, command := range commands {
		if command.order != 0 {
			orderedCommands = append(orderedCommands, command)
		} else {
			remainingCommands = append(remainingCommands, command.command)
		}
	}

	slices.SortFunc(orderedCommands, func(a, b command) int {
		return cmp.Compare(a.order, b.order)
	})

	slices.SortFunc(remainingCommands, func(a, b history.Command) int {
		return cmp.Compare(a.Index, b.Index)
	})

	var commandList []history.Command
	for _, command := range orderedCommands {
		commandList = append(commandList, command.command)
	}

	return append(commandList, remainingCommands...)
}
