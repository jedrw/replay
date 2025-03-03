package tui

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/jedrw/replay/internal/command"
	"github.com/rivo/tview"
)

const (
	ORDERCOLUMN   = iota
	COMMANDCOLUMN = iota
)

func newCommandSelect(replayTui *replayTui) *tview.Table {
	commandSelect := tview.NewTable()
	commandSelect.SetTitle(" Select ").SetTitleAlign(tview.AlignLeft)
	commandSelect.SetBorder(true).SetBorderPadding(0, 0, 1, 1)
	commandSelect.SetBackgroundColor(tcell.ColorDefault)
	commandSelect.SetFixed(1, 0)
	commandSelect.SetSelectable(true, false).SetSeparator(tview.Borders.Vertical)

	commandSelect.SetSelectedFunc(func(row, col int) {
		commandCell := commandSelect.GetCell(row, COMMANDCOLUMN)
		orderCell := commandSelect.GetCell(row, ORDERCOLUMN)
		replayTui.toggleSelected(commandCell, orderCell)
		replayTui.updatePreview()
	})

	commandSelect.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
	)

	return commandSelect
}

func fKeyToNumber(event *tcell.EventKey) int {
	return int(event.Key() - 278)
}

func isSelected(commandCell *tview.TableCell) bool {
	return commandCell.Style == selectedStyle
}

func (replayTui *replayTui) commandInSelectedList(command command.Command) (bool, *tuiCommand) {
	for _, selectedCommand := range replayTui.Selected {
		if selectedCommand.Command.Index == command.Index {
			return true, &selectedCommand
		}
	}

	return false, nil
}

func (replayTui *replayTui) selectCommand(commandCell, orderCell *tview.TableCell) {
	commandCell.SetStyle(selectedStyle)
	replayTui.Selected = append(
		replayTui.Selected,
		tuiCommand{
			Order:   orderCell.Reference.(int),
			Command: commandCell.Reference.(command.Command),
		},
	)
}

func (replayTui *replayTui) deselectCommand(commandCell, orderCell *tview.TableCell) {
	commandCell.SetStyle(unselectedStyle)
	orderCell.SetText("").SetReference(0)
	for i, c := range replayTui.Selected {
		if c.Command.Index == commandCell.Reference.(command.Command).Index {
			replayTui.Selected = append(replayTui.Selected[:i], replayTui.Selected[i+1:]...)
		}
	}
}

func (replayTui *replayTui) toggleSelected(commandCell *tview.TableCell, orderCell *tview.TableCell) {
	if isSelected(commandCell) {
		replayTui.deselectCommand(commandCell, orderCell)
	} else {
		replayTui.selectCommand(commandCell, orderCell)
	}
}

func (replayTui *replayTui) populateCommandTable(commands []command.Command) {
	for column, header := range []string{"Order", "Command"} {
		headerCell := tview.NewTableCell(header).SetSelectable(false)
		headerCell.SetTextColor(altColour)
		headerCell.SetStyle(tcell.Style.Attributes(headerCell.Style, tcell.AttrBold))
		replayTui.commandSelect.SetCell(0, column, headerCell)
	}

	for row, command := range commands {
		commandCell := tview.NewTableCell(command.Command)
		commandCell.SetReference(command)
		replayTui.commandSelect.SetCell(row+1, COMMANDCOLUMN, commandCell)

		orderCell := tview.NewTableCell("").SetReference(0)
		selected, selectedCommand := replayTui.commandInSelectedList(command)
		if selected {
			commandCell.SetStyle(selectedStyle)
			if selectedCommand.Order != 0 {
				orderCell.SetText(strconv.Itoa(selectedCommand.Order))
			}
			orderCell.SetReference(selectedCommand.Order)
		}

		orderCell.SetSelectable(false)
		replayTui.commandSelect.SetCell(row+1, ORDERCOLUMN, orderCell)
	}

	replayTui.commandSelect.Select(replayTui.commandSelect.GetRowCount()-1, COMMANDCOLUMN)
}
