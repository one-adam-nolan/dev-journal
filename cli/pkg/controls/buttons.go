package controls

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var btnStyle tcell.Style = tcell.StyleDefault.
	Foreground(tcell.Color164).
	Background(tcell.Color234)
	// Attributes(tcellwhicih .AttrUnderline)

func GetButton(txt string, action func()) *tview.Button {
	btn := tview.NewButton(txt).
		SetActivatedStyle(tcell.StyleDefault.
			Background(tcell.Color236).
			Attributes(tcell.AttrUnderline)).
		SetSelectedFunc(action)

	btn.SetStyle(btnStyle)

	btn.SetBorder(true)

	return btn
}

type TviewTabControl struct {
	Index    *int
	Controls []interface{}
	HostView *tview.Grid
	App      *tview.Application
}

func NewTabControl(
	controls []interface{},
	host *tview.Grid,
	app *tview.Application) *TviewTabControl {
	index := 0
	host.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if index == (len(controls) - 1) {
				index = 0
			} else {
				index++
			}

			app.SetFocus(controls[index].(tview.Primitive))
		}

		return event
	})
	return &TviewTabControl{
		Index:    &index,
		Controls: controls,
		HostView: host,
		App:      app,
	}
}
