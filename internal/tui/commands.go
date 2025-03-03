package tui

import (
	"cmp"
	"slices"

	"github.com/jedrw/replay/internal/command"
)

func sortCommands(commands []tuiCommand) []command.Command {
	var orderedCommands []tuiCommand
	var remainingCommands []command.Command
	for _, command := range commands {
		if command.Order != 0 {
			orderedCommands = append(orderedCommands, command)
		} else {
			remainingCommands = append(remainingCommands, command.Command)
		}
	}

	slices.SortFunc(orderedCommands, func(a, b tuiCommand) int {
		return cmp.Compare(a.Order, b.Order)
	})

	slices.SortFunc(remainingCommands, func(a, b command.Command) int {
		return cmp.Compare(a.Index, b.Index)
	})

	var commandList []command.Command
	for _, command := range orderedCommands {
		commandList = append(commandList, command.Command)
	}

	return append(commandList, remainingCommands...)
}
