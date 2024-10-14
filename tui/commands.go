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
		if command.Order != 0 {
			orderedCommands = append(orderedCommands, command)
		} else {
			remainingCommands = append(remainingCommands, command.Command)
		}
	}

	slices.SortFunc(orderedCommands, func(a, b command) int {
		return cmp.Compare(a.Order, b.Order)
	})

	slices.SortFunc(remainingCommands, func(a, b history.Command) int {
		return cmp.Compare(a.Index, b.Index)
	})

	var commandList []history.Command
	for _, command := range orderedCommands {
		commandList = append(commandList, command.Command)
	}

	return append(commandList, remainingCommands...)
}
