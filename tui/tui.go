package tui

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/lupinelab/replay/history"
	"github.com/rivo/tview"
)

type command struct {
	Order   int
	Command history.Command
}

type replayTui struct {
	app    *tview.Application
	layout *tview.Flex

	footer        *tview.TextView
	search        *tview.InputField
	commandSelect *tview.Table
	preview       *tview.TextView

	history  history.CommandHistory
	Selected []command
}

func NewReplayTui() replayTui {
	return replayTui{}
}

func (replayTui *replayTui) Run() error {
	replayTui.app = tview.NewApplication()
	replayTui.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	replayTui.layout.SetBorder(true).SetBorderColor(tcell.ColorBlack).SetTitle(" REPLAY ")
	replayTui.layout.SetBackgroundColor(tcell.ColorDefault)

	replayTui.search = newSearch(replayTui)
	replayTui.layout.AddItem(replayTui.search, 3, 0, false)

	replayTui.commandSelect = newCommandSelect(replayTui)
	replayTui.layout.AddItem(replayTui.commandSelect, 0, 1, false)

	replayTui.preview = newPreview()
	replayTui.layout.AddItem(replayTui.preview, 0, 1, false)

	replayTui.footer = NewFooter()
	replayTui.layout.AddItem(replayTui.footer, 1, 0, false)

	var err error
	replayTui.history, err = history.GetHistory()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	replayTui.populateCommandTable(replayTui.history)
	replayTui.app.SetInputCapture(replayTui.inputHandler())
	replayTui.app.SetRoot(replayTui.layout, true)

	return replayTui.app.Run()
}

func isFKey(event *tcell.EventKey) bool {
	if 279 <= event.Key() && event.Key() <= 288 {
		return true
	}

	return false
}

func isCommandSelectInput(event *tcell.EventKey) bool {
	if (event.Modifiers() == tcell.ModAlt && event.Key() == tcell.KeyEnter) ||
		(event.Key() == tcell.KeyCtrlR) ||
		(event.Key() == tcell.KeyEnter || event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown) {
		return true
	} else if isFKey(event) {
		return true
	} else {
		return false
	}
}

func (replayTui *replayTui) inputHandler() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyCtrlC {
			replayTui.app.Stop()
			return nil
		} else if isCommandSelectInput(event) {
			handleCommandSelectInput := replayTui.commandSelect.InputHandler()
			handleCommandSelectInput(event, nil)
			return nil
		} else {
			handleSearchInput := replayTui.search.InputHandler()
			handleSearchInput(event, nil)
			return nil
		}
	}
}
