package tui_test

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/jedrw/replay/history"
	"github.com/jedrw/replay/tui"
	"github.com/rivo/tview"
)

func TestEventRuneToNumberKey(t *testing.T) {
	eventKey1 := tcell.NewEventKey(tcell.KeyF1, 0, tcell.ModNone)
	eventKey2 := tcell.NewEventKey(tcell.KeyF9, 0, tcell.ModNone)

	if number := tui.FKeyToNumber(eventKey1); number != 1 {
		t.Errorf("Expected: 1, got: %d", number)
	}
	if number := tui.FKeyToNumber(eventKey2); number != 9 {
		t.Errorf("Expected: 9, got: %d", number)
	}
}

func TestIsSelectedReturnsTrue(t *testing.T) {
	commandCell := tview.NewTableCell("").SetStyle(tui.SelectedStyle)
	if isSelected := tui.IsSelected(commandCell); !isSelected {
		t.Errorf("Expected: true, Got: %t", isSelected)
	}
}

func TestIsSelectedReturnsFalse(t *testing.T) {
	commandCell := tview.NewTableCell("").SetStyle(tui.UnselectedStyle)
	if isSelected := tui.IsSelected(commandCell); isSelected {
		t.Errorf("Expected: false, Got: %t", isSelected)
	}
}

func TestCommandInSelectedListRetunsCommand(t *testing.T) {
	replayTui := tui.ReplayTui{}
	cmd := history.Command{
		Index:   0,
		Command: "echo foo",
	}

	selectedCommand := tui.Command{
		Order:   0,
		Command: cmd,
	}

	replayTui.Selected = append(
		replayTui.Selected,
		selectedCommand,
		tui.Command{
			Order: 0,
			Command: history.Command{
				Index:   1,
				Command: "echo bar",
			},
		},
	)

	found, foundCommand := replayTui.CommandInSelectedList(cmd)
	if !found {
		t.Errorf("Expected to find command but returned: %t", found)
	}

	if foundCommand.Command.Index != cmd.Index {
		t.Errorf("Expected:%d, Got:%d", cmd.Index, foundCommand.Command.Index)
	}

	if foundCommand.Command.Command != cmd.Command {
		t.Errorf("Expected:%s, Got:%s", cmd.Command, foundCommand.Command.Command)
	}
}

func TestCommandInSelectedListReturnsNil(t *testing.T) {
	replayTui := tui.ReplayTui{}
	cmd := history.Command{
		Index:   0,
		Command: "echo foo",
	}

	replayTui.Selected = append(
		replayTui.Selected,
		tui.Command{
			Order: 0,
			Command: history.Command{
				Index:   1,
				Command: "echo bar",
			},
		},
	)

	found, foundCommand := replayTui.CommandInSelectedList(cmd)
	if found {
		t.Errorf("Expected not to find command but returned: %t", found)
	}

	if foundCommand != nil {
		t.Errorf("Expected not to find command but returned: %v", foundCommand)
	}
}

func TestSelectCommand(t *testing.T) {
	replayTui := tui.ReplayTui{}
	cmd := tui.Command{
		Order: 0,
		Command: history.Command{
			Index:   0,
			Command: "echo foo",
		},
	}

	commandCell := tview.NewTableCell("echo foo").SetReference(cmd.Command)
	orderCell := tview.NewTableCell("").SetReference(0)

	replayTui.SelectCommand(commandCell, orderCell)

	if commandCell.Style != tui.SelectedStyle {
		t.Errorf("\nExpected:\n%+v\nGot:\n%+v", tui.SelectedStyle, commandCell.Style)
	}

	var selectedCmd tui.Command
	for _, command := range replayTui.Selected {
		if command.Command.Index == cmd.Command.Index {
			selectedCmd = command
		}
	}

	if selectedCmd != cmd {
		t.Errorf("\nExpected:\n%+v\nGot:\n%+v", cmd, selectedCmd)
	}
}

func TestDeselectCommand(t *testing.T) {
	replayTui := tui.ReplayTui{}
	cmd := tui.Command{
		Order: 1,
		Command: history.Command{
			Index:   0,
			Command: "echo foo",
		},
	}

	replayTui.Selected = append(replayTui.Selected, cmd)

	commandCell := tview.NewTableCell("echo foo").SetReference(cmd.Command)
	orderCell := tview.NewTableCell("1").SetReference(1)

	replayTui.DeselectCommand(commandCell, orderCell)

	if commandCell.Style != tui.UnselectedStyle {
		t.Errorf("\nExpected:\n%+v\nGot:\n%+v", tui.UnselectedStyle, commandCell.Style)
	}

	for _, command := range replayTui.Selected {
		if command.Command.Index == cmd.Command.Index {
			t.Errorf("Expected selected list not to contain:\n%+v", cmd)
		}
	}
}
