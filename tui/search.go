package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/jedrw/replay/command"
	"github.com/rivo/tview"
)

func newSearch(replayTui *replayTui) *tview.InputField {
	search := tview.NewInputField()
	search.SetTitle(" Search ").SetTitleAlign(tview.AlignLeft)
	search.SetBorder(true).SetBorderPadding(0, 0, 1, 1)
	search.SetBackgroundColor(tcell.ColorDefault)
	search.SetFieldStyle(tcell.StyleDefault)

	search.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return event
	})

	search.SetChangedFunc(func(searchText string) {
		if searchText == "" {
			replayTui.commandSelect.Clear()
			replayTui.populateCommandTable(replayTui.shellHistory)
			replayTui.commandSelect.Select(replayTui.commandSelect.GetRowCount()-1, COMMANDCOLUMN)
			return
		}

		matches := replayTui.searchHistory(searchText)

		replayTui.commandSelect.Clear()
		replayTui.populateCommandTable(matches)
		replayTui.commandSelect.Select(replayTui.commandSelect.GetRowCount()-1, COMMANDCOLUMN)
	})

	return search
}

func (replayTui *replayTui) searchHistory(searchText string) []command.Command {
	var matches []command.Command
	for _, command := range replayTui.shellHistory {
		if strings.Contains(command.Command, searchText) {
			matches = append(matches, command)
		}
	}

	return matches
}
