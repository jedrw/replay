package tui

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/jedrw/replay/command"
	"github.com/jedrw/replay/history"
	"github.com/jedrw/replay/replay"
	"github.com/rivo/tview"
)

const commandPage = "command"
const historyPage = "history"

type tuiCommand struct {
	Order   int
	Command command.Command
}

type replayTui struct {
	app    *tview.Application
	pages  *tview.Pages
	layout *tview.Flex

	footer        *tview.TextView
	search        *tview.InputField
	commandSelect *tview.Table
	preview       *tview.TextView
	history       *tview.Flex
	historyPages  *tview.Pages

	shellHistory      command.ShellHistory
	replayHistory     history.ReplayHistory
	replayHistoryPath string
	Selected          []tuiCommand
}

func NewReplayTui() replayTui {
	return replayTui{}
}

func (replayTui *replayTui) Run() error {
	replayTui.app = tview.NewApplication()
	replayTui.pages = tview.NewPages()

	replayTui.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	replayTui.layout.SetBorder(true).SetBorderColor(tcell.ColorBlack).SetTitle(" REPLAY ")
	replayTui.layout.SetBackgroundColor(tcell.ColorDefault)

	replayTui.search = newSearch(replayTui)
	replayTui.layout.AddItem(replayTui.search, 3, 0, false)

	replayTui.commandSelect = newCommandSelect(replayTui)
	replayTui.layout.AddItem(replayTui.commandSelect, 0, 1, false)

	replayTui.preview = newPreview()
	replayTui.layout.AddItem(replayTui.preview, 0, 1, false)

	replayTui.footer = newFooter()
	replayTui.layout.AddItem(replayTui.footer, 1, 0, false)

	replayTui.pages.AddPage(commandPage, replayTui.layout, true, true)

	var err error
	shellHistoryPath, err := command.ShellHistoryPath()
	if err != nil {
		return err
	}

	replayTui.shellHistory, err = command.GetShellHistory(shellHistoryPath)
	if err != nil {
		return err
	}

	replayTui.populateCommandTable(replayTui.shellHistory)
	replayTui.replayHistoryPath, err = history.ReplayHistoryPath()
	if err != nil {
		return err
	}

	replayTui.replayHistory, err = history.GetReplayHistory(replayTui.replayHistoryPath)
	if err != nil {
		return err
	}

	replayTui.history = newHistory(replayTui)
	replayTui.pages.AddPage(historyPage, replayTui.history, true, false)

	replayTui.pages.SetInputCapture(replayTui.inputHandler())
	replayTui.app.SetRoot(replayTui.pages, true)
	replayTui.app.SetFocus(replayTui.pages)

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
		(event.Key() == tcell.KeyCtrlR || event.Key() == tcell.KeyCtrlH || event.Key() == tcell.KeyEnter || event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown) {
		return true
	} else if isFKey(event) {
		return true
	} else {
		return false
	}
}

func (replayTui *replayTui) inputHandler() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		pageName, _ := replayTui.pages.GetFrontPage()
		switch pageName {
		case historyPage:
			if event.Key() == tcell.KeyEsc {
				replayTui.pages.HidePage(historyPage)
				replayTui.footer.SetText(commandSelectControls)
			} else if event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown {
				pageString, _ := replayTui.historyPages.GetFrontPage()
				pageNum, err := strconv.Atoi(pageString)
				if err != nil {
					return nil
				}

				numHistory := len(replayTui.replayHistory)

				direction := 0
				if event.Key() == tcell.KeyUp {
					if pageNum < numHistory-1 {
						direction = 1
					}
				} else {
					if pageNum > 0 {
						direction = -1
					}
				}

				replayTui.historyPages.SwitchToPage(strconv.Itoa(int(math.Abs(float64((pageNum + direction) % replayTui.historyPages.GetPageCount())))))
			} else if event.Key() == tcell.KeyCtrlR {
				pageName, _ := replayTui.historyPages.GetFrontPage()
				pageNum, err := strconv.Atoi(pageName)
				if err != nil {
					replayTui.app.Stop()
					fmt.Println(err)
					return nil
				}

				r := replayTui.replayHistory[pageNum]
				var commands []command.Command
				for _, c := range r {
					commands = append(commands, command.Command{
						Command: c,
					})
				}

				replayTui.app.Stop()
				replay.Replay(commands)
				newReplay := history.NewReplayFromCommands(commands)
				replayTui.replayHistory = history.UpdateReplayHistory(newReplay, replayTui.replayHistory)
				history.WriteReplayHistory(replayTui.replayHistory, replayTui.replayHistoryPath)
			}
		case commandPage:
			if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyCtrlC {
				replayTui.app.Stop()
			} else if isCommandSelectInput(event) {
				if (event.Modifiers() == tcell.ModAlt && event.Key() == tcell.KeyEnter) ||
					(event.Key() == tcell.KeyCtrlR) {
					commands := sortCommands(replayTui.Selected)
					replayTui.app.Stop()
					replay.Replay(commands)
					newReplay := history.NewReplayFromCommands(commands)
					newHistory := history.UpdateReplayHistory(newReplay, replayTui.replayHistory)
					history.WriteReplayHistory(newHistory, replayTui.replayHistoryPath)
				} else if event.Key() == tcell.KeyEnter || event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown {
					commandSelectInputHandler := replayTui.commandSelect.InputHandler()
					commandSelectInputHandler(event, func(p tview.Primitive) {})
				} else if event.Key() == tcell.KeyEsc {
					replayTui.app.Stop()
				} else if event.Key() == tcell.KeyCtrlH {
					replayTui.pages.ShowPage(historyPage)
					replayTui.footer.SetText(historyControls)
				} else if isFKey(event) {
					switch order := fKeyToNumber(event); order {
					case 1, 2, 3, 4, 5, 6, 7, 8, 9:
						// Check for another command with same order
						for i, selectedCommand := range replayTui.Selected {
							if selectedCommand.Order == order {
								// The command might not be in the commandSelect table so if there has
								// been a search it must be removed from the selected list manually
								replayTui.Selected = append(replayTui.Selected[:i], replayTui.Selected[i+1:]...)
								for r := 1; r < replayTui.commandSelect.GetRowCount(); r++ {
									commandCell := replayTui.commandSelect.GetCell(r, COMMANDCOLUMN)
									if commandCell.Reference.(command.Command).Index == selectedCommand.Command.Index {
										orderCell := replayTui.commandSelect.GetCell(r, ORDERCOLUMN)
										orderCell.SetText("").SetReference(0)
										replayTui.deselectCommand(commandCell, orderCell)
										break
									}
								}
							}
						}

						row, _ := replayTui.commandSelect.GetSelection()
						commandCell := replayTui.commandSelect.GetCell(row, COMMANDCOLUMN)
						orderCell := replayTui.commandSelect.GetCell(row, ORDERCOLUMN)
						orderCell.SetText(strconv.Itoa(order)).SetReference(order)

						// Set order on already selected command
						if isSelected(commandCell) {
							for i, c := range replayTui.Selected {
								if c.Command.Index == commandCell.Reference.(command.Command).Index {
									replayTui.Selected[i].Order = order
									break
								}
							}
						} else {
							replayTui.selectCommand(commandCell, orderCell)
						}

						replayTui.updatePreview()
					}

				}
			} else {
				handleSearchInput := replayTui.search.InputHandler()
				handleSearchInput(event, nil)
			}
		}

		return nil
	}
}
