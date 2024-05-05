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
	commandSelect.SetBorder(true).SetTitle(" Select ").SetBorderPadding(0, 0, 1, 1)
	commandSelect.SetBackgroundColor(tcell.ColorDefault)
	commandSelect.SetFixed(1, 0)
	commandSelect.SetSelectable(true, false).SetSeparator(tview.Borders.Vertical)

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
			return nil
		} else if event.Key() == tcell.KeyEnter || event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown || event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyRight {
			return event
		} else if event.Key() == tcell.KeyEsc {
			replayTui.app.Stop()
			return nil
		} else if number := eventRuneToNumberKey(event); 1 <= number && number <= 9 {
			switch eventRuneToNumberKey(event) {
			case 1, 2, 3, 4, 5, 6, 7, 8, 9:
				order := eventRuneToNumberKey(event)
				// Check for another command with same order
				for i, selectedCommand := range replayTui.selected {
					if selectedCommand.order == order {
						// The command might not be in the commandSelect table so if there has
						// been a search so it must be removed from the selected list manually
						replayTui.selected = append(replayTui.selected[:i], replayTui.selected[i+1:]...)
						for r := 1; r < commandSelect.GetRowCount(); r++ {
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
			return nil
		} else {
			handleSearchInput := replayTui.search.InputHandler()
			handleSearchInput(event, nil)
			return nil
		}
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

func (replayTui *replayTui) commandInSelectedList(command history.Command) (bool, *command) {
	for _, selectCommand := range replayTui.selected {
		if selectCommand.command.Index == command.Index {
			return true, &selectCommand
		}
	}

	return false, nil
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

func (replayTui *replayTui) populateCommandTable(commands []history.Command) {
	for column, header := range []string{"Index", "Command", "Order"} {
		headerCell := tview.NewTableCell(header).SetSelectable(false)
		headerCell.SetTextColor(tcell.ColorSteelBlue)
		headerCell.SetStyle(tcell.Style.Attributes(headerCell.Style, tcell.AttrBold))
		replayTui.commandSelect.SetCell(0, column, headerCell)
	}

	for row, command := range commands {
		indexCell := tview.NewTableCell(fmt.Sprint(command.Index))
		indexCell.SetSelectable(false)
		replayTui.commandSelect.SetCell(row+1, 0, indexCell)

		commandCell := tview.NewTableCell(fmt.Sprint(command.Command))
		commandCell.SetReference(command)
		replayTui.commandSelect.SetCell(row+1, 1, commandCell)

		orderCell := tview.NewTableCell("").SetReference(0)
		selected, selectedCommand := replayTui.commandInSelectedList(command)
		if selected {
			commandCell.SetStyle(selectedStyle)
			if selectedCommand.order != 0 {
				orderCell.SetText(fmt.Sprint(selectedCommand.order))
			}
			orderCell.SetReference(selectedCommand.order)
		}

		orderCell.SetSelectable(false)
		replayTui.commandSelect.SetCell(row+1, 2, orderCell)
	}
}
