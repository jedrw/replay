package tui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/lupinelab/replay/history"
	"github.com/rivo/tview"
)

func TestEventRuneToNumberKey(t *testing.T) {
	eventKey1 := tcell.NewEventKey(tcell.KeyF1, 0, tcell.ModNone)
	eventKey2 := tcell.NewEventKey(tcell.KeyF9, 0, tcell.ModNone)

	if number := fKeyToNumber(eventKey1); number != 1 {
		t.Errorf("Expected: 1, got: %d", number)
	}
	if number := fKeyToNumber(eventKey2); number != 9 {
		t.Errorf("Expected: 9, got: %d", number)
	}
}

func TestIsSelectedReturnsTrue(t *testing.T) {
	commandCell := tview.NewTableCell("").SetStyle(selectedStyle)
	if isSelected := isSelected(commandCell); !isSelected {
		t.Errorf("Expected: true, Got: %t", isSelected)
	}
}

func TestIsSelectedReturnsFalse(t *testing.T) {
	commandCell := tview.NewTableCell("").SetStyle(unselectedStyle)
	if isSelected := isSelected(commandCell); isSelected {
		t.Errorf("Expected: false, Got: %t", isSelected)
	}
}

func TestCommandInSelectedListRetunsCommand(t *testing.T) {
	replayTui := replayTui{}
	cmd := history.Command{
		Index:   0,
		Command: "echo foo",
	}

	selectedCommand := command{
		order:   0,
		command: cmd,
	}

	replayTui.selected = append(
		replayTui.selected,
		selectedCommand,
		command{
			order: 0,
			command: history.Command{
				Index:   1,
				Command: "echo bar",
			},
		},
	)

	found, foundCommand := replayTui.commandInSelectedList(cmd)
	if !found {
		t.Errorf("Expected to find command but returned: %t", found)
	}

	if foundCommand.command.Index != cmd.Index {
		t.Errorf("Expected:%d, Got:%d", cmd.Index, foundCommand.command.Index)
	}

	if foundCommand.command.Command != cmd.Command {
		t.Errorf("Expected:%s, Got:%s", cmd.Command, foundCommand.command.Command)
	}
}

func TestCommandInSelectedListReturnsNil(t *testing.T) {
	replayTui := replayTui{}
	cmd := history.Command{
		Index:   0,
		Command: "echo foo",
	}

	replayTui.selected = append(
		replayTui.selected,
		command{
			order: 0,
			command: history.Command{
				Index:   1,
				Command: "echo bar",
			},
		},
	)

	found, foundCommand := replayTui.commandInSelectedList(cmd)
	if found {
		t.Errorf("Expected not to find command but returned: %t", found)
	}

	if foundCommand != nil {
		t.Errorf("Expected not to find command but returned: %v", foundCommand)
	}
}

func TestSelectCommand(t *testing.T) {
	replayTui := replayTui{}
	cmd := command{
		order: 0,
		command: history.Command{
			Index:   0,
			Command: "echo foo",
		},
	}

	commandCell := tview.NewTableCell("echo foo").SetReference(cmd.command)
	orderCell := tview.NewTableCell("").SetReference(0)

	replayTui.selectCommand(commandCell, orderCell)

	if commandCell.Style != selectedStyle {
		t.Errorf("\nExpected:\n%+v\nGot:\n%+v", selectedStyle, commandCell.Style)
	}

	var selectedCmd command
	for _, command := range replayTui.selected {
		if command.command.Index == cmd.command.Index {
			selectedCmd = command
		}
	}

	if selectedCmd != cmd {
		t.Errorf("\nExpected:\n%+v\nGot:\n%+v", cmd, selectedCmd)
	}
}

func TestDeselectCommand(t *testing.T) {
	replayTui := replayTui{}
	cmd := command{
		order: 1,
		command: history.Command{
			Index:   0,
			Command: "echo foo",
		},
	}

	replayTui.selected = append(replayTui.selected, cmd)

	commandCell := tview.NewTableCell("echo foo").SetReference(cmd.command)
	orderCell := tview.NewTableCell("1").SetReference(1)

	replayTui.deselectCommand(commandCell, orderCell)

	if commandCell.Style != unselectedStyle {
		t.Errorf("\nExpected:\n%+v\nGot:\n%+v", unselectedStyle, commandCell.Style)
	}

	for _, command := range replayTui.selected {
		if command.command.Index == cmd.command.Index {
			t.Errorf("Expected selected list not to contain:\n%+v", cmd)
		}
	}
}
