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

	header        *tview.TextView
	search        *tview.InputField
	commandSelect *tview.Table
	preview       *tview.TextView

	history  history.CommandHistory
	selected []command
}

func NewReplayTui() replayTui {
	return replayTui{}
}

func (replayTui *replayTui) Run() error {
	replayTui.app = tview.NewApplication()
	replayTui.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	replayTui.layout.SetBorder(true).SetBorderColor(tcell.ColorBlack).SetTitle(" REPLAY ")
	replayTui.layout.SetBackgroundColor(tcell.ColorDefault)

	replayTui.header = NewHeader()
	replayTui.layout.AddItem(replayTui.header, 1, 0, false)

	replayTui.search = newSearch(replayTui)
	replayTui.layout.AddItem(replayTui.search, 3, 0, false)

	replayTui.commandSelect = newCommandSelect(replayTui)
	replayTui.layout.AddItem(replayTui.commandSelect, 0, 1, true)

	replayTui.preview = newPreview()
	replayTui.layout.AddItem(replayTui.preview, 0, 1, false)

	var err error
	replayTui.history, err = history.GetHistory()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	replayTui.populateCommandTable(replayTui.history)
	replayTui.app.SetRoot(replayTui.layout, true)

	return replayTui.app.Run()
}
