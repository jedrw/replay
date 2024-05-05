package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/lupinelab/replay/history"
	"github.com/lupinelab/replay/replay.go"
	"github.com/rivo/tview"
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
		commandCell := commandSelect.GetCell(row, col)
		orderCell := commandSelect.GetCell(row, col+1)
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

		switch event.Rune() - 48 {
		case 1, 2, 3, 4, 5, 6, 7, 8, 9:
			order := int(event.Rune() - 48)

			// Check for another command with same order
			for r := 1; r < commandSelect.GetRowCount(); r++ {
				commandCell := commandSelect.GetCell(r, 1)
				orderCell := commandSelect.GetCell(r, 2)
				if orderCell.Reference.(int) == order {
					orderCell.SetText("").SetReference(0)
					replayTui.toggleSelected(commandCell, orderCell)
				}
			}

			row, _ := commandSelect.GetSelection()
			commandCell := commandSelect.GetCell(row, 1)
			orderCell := commandSelect.GetCell(row, 2)

			// Set order on selected command
			if isSelected(commandCell) {
				for i, command := range replayTui.selected {
					if command.command.Index == commandCell.Reference.(history.Command).Index {
						replayTui.selected[i].order = order
						orderCell.SetText(fmt.Sprint(order)).SetReference(order)
						updatePreview(replayTui)
						return event
					}
				}
			}

			orderCell.SetText(fmt.Sprint(order)).SetReference(order)
			selectCommand(commandCell, orderCell, replayTui)
		}

		updatePreview(replayTui)
		return event
	})

	commandSelect.SetFocusFunc(func() {
		commandSelect.Select(commandSelect.GetRowCount()-1, 0)
	})

	return commandSelect
}

func isSelected(commandCell *tview.TableCell) bool {
	return commandCell.Style == selectedStyle
}

func selectCommand(commandCell, orderCell *tview.TableCell, replayTui *replayTui) {
	commandCell.SetStyle(selectedStyle)
	replayTui.selected = append(
		replayTui.selected,
		command{
			order:   orderCell.Reference.(int),
			command: commandCell.Reference.(history.Command),
		},
	)
}

func deselectCommand(commandCell, orderCell *tview.TableCell, replayTui *replayTui) {
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
		deselectCommand(commandCell, orderCell, replayTui)
	} else {
		selectCommand(commandCell, orderCell, replayTui)
	}
}
