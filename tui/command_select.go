package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/lupinelab/replay/history"
	"github.com/lupinelab/replay/replay"
	"github.com/rivo/tview"
)

const (
	INDEXCOLUMN   = iota
	COMMANDCOLUMN = iota
	ORDERCOLUMN   = iota
)

func newCommandSelect(replayTui *replayTui) *tview.Table {
	commandSelect := tview.NewTable()
	commandSelect.SetBorder(true)
	commandSelect.SetBackgroundColor(tcell.ColorDefault)
	commandSelect.SetFixed(1, 0)
	commandSelect.SetSelectable(true, false).SetSeparator(tview.Borders.Vertical)

	for column, header := range []string{"Index", "Command", "Order"} {
		headerCell := tview.NewTableCell(header).SetSelectable(false)
		headerCell.SetTextColor(tcell.ColorGrey)
		headerCell.SetStyle(tcell.Style.Attributes(headerCell.Style, tcell.AttrBold))
		commandSelect.SetCell(0, column, headerCell)
	}

	commandSelect.SetSelectedFunc(func(row, col int) {
		commandCell := commandSelect.GetCell(row, COMMANDCOLUMN)
		orderCell := commandSelect.GetCell(row, ORDERCOLUMN)
		replayTui.toggleSelected(commandCell, orderCell)
		updatePreview(replayTui)
	})

	commandSelect.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers() == tcell.ModAlt && event.Key() == tcell.KeyEnter {
			commands := sortCommands(replayTui.selected)
			replay := replay.NewReplay(commands)
			replayTui.app.Stop()
			replay.Run()
		}

		switch eventRuneToNumberKey(event) {
		case 1, 2, 3, 4, 5, 6, 7, 8, 9:
			order := eventRuneToNumberKey(event)
			// Check for another command with same order
			for _, selectedCommand := range replayTui.selected {
				if selectedCommand.order == order {
					for r := 1; r < commandSelect.GetRowCount()-1; r++ {
						commandCell := commandSelect.GetCell(r, COMMANDCOLUMN)
						if commandCell.Reference.(history.Command).Index == selectedCommand.command.Index {
							orderCell := commandSelect.GetCell(r, ORDERCOLUMN)
							orderCell.SetText("").SetReference(0)
							replayTui.deselectCommand(commandCell, orderCell)
							break
						}
					}
				}
			}

			row, _ := commandSelect.GetSelection()
			commandCell := commandSelect.GetCell(row, COMMANDCOLUMN)
			orderCell := commandSelect.GetCell(row, ORDERCOLUMN)
			orderCell.SetText(fmt.Sprint(order)).SetReference(order)

			// Set order on already selected command
			if isSelected(commandCell) {
				for i, command := range replayTui.selected {
					if command.command.Index == commandCell.Reference.(history.Command).Index {
						replayTui.selected[i].order = order
						break
					}
				}
			} else {
				replayTui.selectCommand(commandCell, orderCell)
			}

			updatePreview(replayTui)
		}
		return event
	})

	commandSelect.SetFocusFunc(func() {
		commandSelect.Select(commandSelect.GetRowCount()-1, COMMANDCOLUMN)
	})

	return commandSelect
}

func eventRuneToNumberKey(event *tcell.EventKey) int {
	return int(event.Rune() - '0')
}

func isSelected(commandCell *tview.TableCell) bool {
	return commandCell.Style == selectedStyle
}

func (replayTui *replayTui) selectCommand(commandCell, orderCell *tview.TableCell) {
	commandCell.SetStyle(selectedStyle)
	replayTui.selected = append(
		replayTui.selected,
		command{
			order:   orderCell.Reference.(int),
			command: commandCell.Reference.(history.Command),
		},
	)
}

func (replayTui *replayTui) deselectCommand(commandCell, orderCell *tview.TableCell) {
	commandCell.SetStyle(unselectedStyle)
	orderCell.SetText("").SetReference(0)
	for i, command := range replayTui.selected {
		if command.command.Index == commandCell.Reference.(history.Command).Index {
			replayTui.selected = append(replayTui.selected[:i], replayTui.selected[i+1:]...)
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
