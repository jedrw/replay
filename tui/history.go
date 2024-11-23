package tui

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func newHistory(replayTui *replayTui) *tview.Flex {
	replayTui.historyPages = tview.NewPages()
	history := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(replayTui.historyPages, 0, 10, false).
				AddItem(nil, 0, 1, false), 0, 10, false).
		AddItem(nil, 0, 1, false)

	for i, replay := range replayTui.replayHistory {
		r := tview.NewTextView()
		r.SetBackgroundColor(tcell.ColorDefault)

		r.SetTitle(fmt.Sprintf(" Replay History %d/%d ", i+1, len(replayTui.replayHistory))).SetTitleAlign(tview.AlignLeft).
			SetBorder(true).
			SetBorderPadding(0, 0, 1, 1).
			SetRect(2, 2, 4, 5)

		var text string
		for _, command := range replay {
			text += command + "\n"
		}

		r.SetText(text)
		replayTui.historyPages.AddPage(
			strconv.Itoa(i),
			r,
			true,
			i == 0,
		)
	}

	if len(replayTui.replayHistory) == 0 {
		r := tview.NewTextView()
		r.SetBackgroundColor(tcell.ColorDefault)

		r.SetTitle(" Replay History ").SetTitleAlign(tview.AlignLeft).
			SetBorder(true).
			SetBorderPadding(0, 0, 1, 1).
			SetRect(2, 2, 4, 5)

		r.SetText("empty")
		replayTui.historyPages.AddPage(
			"0",
			r,
			true,
			true,
		)
	}

	return history
}
