package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func newPreview() *tview.TextView {
	preview := tview.NewTextView()
	preview.SetTitle(" Preview ").SetTitleAlign(tview.AlignLeft)
	preview.SetBorder(true).SetBorderPadding(0, 0, 1, 1)
	preview.SetBackgroundColor(tcell.ColorDefault)

	return preview
}

func (replayTui *replayTui) updatePreview() {
	commands := make([]command, len(replayTui.selected))
	copy(commands, replayTui.selected)

	commandList := sortCommands(commands)
	var previewText string
	for _, command := range commandList {
		previewText += command.Command + "\n"
	}

	replayTui.preview.SetText(previewText)
}
