package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TextModalWithQandEscLowerBar struct {
	Content string
}

func (d *TextModalWithQandEscLowerBar) Display() error {
	app := tview.NewApplication()

	// Create a text view to display the content
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetWordWrap(true).
		SetText(d.Content).
		SetDoneFunc(func(key tcell.Key) {
			app.Stop()
		})

	// Set the key event handler to close the window on "Esc" or "q" key press
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Rune() == 'q' {
			app.Stop()
		}
		return event
	})

	statusBar := tview.NewTextView()

	// Set the status bar content
	statusBar.SetText("(q) or esc to quit")
	statusBar.SetTextColor(tcell.ColorOrangeRed)
	statusBar.SetBackgroundColor(tcell.ColorGray)

	// Create a grid layout to hold the text view and status bar
	grid := tview.NewGrid().
		SetRows(-1, 1).
		SetColumns(-1).
		AddItem(textView, 0, 0, 1, 1, 0, 0, true).
		AddItem(statusBar, 1, 0, 1, 1, 0, 0, false)

	// Set the grid layout as the root of the application
	app.SetRoot(grid, true)

	return app.Run()

}
