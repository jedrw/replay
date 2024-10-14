package tui

import (
	"github.com/lupinelab/replay/history"
	"github.com/rivo/tview"
)

type ReplayTui = replayTui

func (r *ReplayTui) CommandInSelectedList(command history.Command) (bool, *command) {
	return r.commandInSelectedList(command)
}

func (r *ReplayTui) SelectCommand(commandCell, orderCell *tview.TableCell) {
	r.selectCommand(commandCell, orderCell)
}

func (r *ReplayTui) DeselectCommand(commandCell, orderCell *tview.TableCell) {
	r.deselectCommand(commandCell, orderCell)
}

type Command = command

var (
	FKeyToNumber    = fKeyToNumber
	SelectedStyle   = selectedStyle
	UnselectedStyle = unselectedStyle
	IsSelected      = isSelected
)
