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

	commandSelect.SetSelectable(true, false).SetSeparator(tview.Borders.Vertical)

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
			for r := 0; r < commandSelect.GetRowCount(); r++ {
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
			for _, command := range replayTui.selected {
				if command.command.Number == commandCell.Reference.(history.Command).Number {
					command.order = order
					orderCell.SetText(fmt.Sprint(order)).SetReference(order)
					updatePreview(replayTui)
					return event
				}
			}
			orderCell.SetText(fmt.Sprint(order)).SetReference(order)
			replayTui.toggleSelected(commandCell, orderCell)
		}

		updatePreview(replayTui)
		return event
	})

	commandSelect.SetFocusFunc(func() {
		commandSelect.Select(commandSelect.GetRowCount()-1, 0)
	})

	return commandSelect
}

func (replayTui *replayTui) toggleSelected(commandCell *tview.TableCell, orderCell *tview.TableCell) {
	if commandCell.Style == selectedStyle {
		commandCell.SetStyle(unselectedStyle)
		orderCell.SetText("").SetReference(0)
		for i, command := range replayTui.selected {
			if command.command.Number == commandCell.Reference.(history.Command).Number {
				replayTui.selected = append(replayTui.selected[:i], replayTui.selected[i+1:]...)
			}
		}
	} else {
		commandCell.SetStyle(selectedStyle)
		replayTui.selected = append(
			replayTui.selected,
			command{
				order:   orderCell.Reference.(int),
				command: commandCell.Reference.(history.Command),
			},
		)
	}
}
