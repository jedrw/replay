package tui

import (
	"github.com/jedrw/replay/command"
	"github.com/rivo/tview"
)

type ReplayTui = replayTui

func (r *ReplayTui) CommandInSelectedList(command command.Command) (bool, *tuiCommand) {
	return r.commandInSelectedList(command)
}

func (r *ReplayTui) SelectCommand(commandCell, orderCell *tview.TableCell) {
	r.selectCommand(commandCell, orderCell)
}

func (r *ReplayTui) DeselectCommand(commandCell, orderCell *tview.TableCell) {
	r.deselectCommand(commandCell, orderCell)
}

type Command = tuiCommand

var (
	FKeyToNumber    = fKeyToNumber
	SelectedStyle   = selectedStyle
	UnselectedStyle = unselectedStyle
	IsSelected      = isSelected
)
