package tui

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/lupinelab/replay/history"
	"github.com/rivo/tview"
)

type command struct {
	order   int
	command history.Command
}

type replayTui struct {
	app    *tview.Application
	layout *tview.Flex

	commandSelect *tview.Table
	preview       *tview.TextView

	selected []command
}

func NewReplayTui() replayTui {
	return replayTui{}
}

func (replayTui *replayTui) Run() error {
	replayTui.app = tview.NewApplication()
	replayTui.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	replayTui.layout.SetBackgroundColor(tcell.ColorDefault)

	replayTui.commandSelect = newCommandSelect(replayTui)
	replayTui.layout.AddItem(replayTui.commandSelect, 0, 1, true)

	replayTui.preview = newPreview()
	replayTui.layout.AddItem(replayTui.preview, 0, 1, false)

	history, err := history.GetHistory()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for row, command := range history {
		numberCell := tview.NewTableCell(fmt.Sprint(command.Number))
		numberCell.SetSelectable(false)
		replayTui.commandSelect.SetCell(row, 0, numberCell)

		commandCell := tview.NewTableCell(fmt.Sprint(command.Command))
		commandCell.SetReference(command)
		replayTui.commandSelect.SetCell(row, 1, commandCell)

		orderCell := tview.NewTableCell("")
		orderCell.SetSelectable(false)
		orderCell.SetReference(0)
		replayTui.commandSelect.SetCell(row, 2, orderCell)
	}

	replayTui.app.SetRoot(replayTui.layout, true)

	return replayTui.app.Run()
}
