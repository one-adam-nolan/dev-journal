package tui

import (
	"dev-journal/directory"
	"dev-journal/pkg/addlogic"
	"dev-journal/pkg/controls"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	ADD_ENTRY  = "input-entry"
	ADD_BULLET = "input-bulet"
)

type TextModalWithQandEscLowerBar struct {
	FilePath       string
	Pages          *tview.Pages
	ContentTxtView *tview.TextView
}

func (d *TextModalWithQandEscLowerBar) UpdateContent() {
	d.ContentTxtView.SetText(d.mustGetTodaysFileText())
}

func (d *TextModalWithQandEscLowerBar) Display() error {
	app := tview.NewApplication()
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:   tcell.Color234,
		PrimaryTextColor:           tcell.Color164,
		BorderColor:                tcell.Color51,
		ContrastSecondaryTextColor: tcell.ColorDarkGreen,
		SecondaryTextColor:         tcell.Color51,
		TertiaryTextColor:          tcell.ColorAliceBlue,
	}

	d.Pages = tview.NewPages()

	d.ContentTxtView = tview.NewTextView().
		SetDynamicColors(true)

	d.UpdateContent()

	statusBar := tview.NewTextView()
	statusBar.SetText("(q) or esc to quit")
	statusBar.SetTextColor(tcell.ColorOrangeRed)
	statusBar.SetBackgroundColor(tcell.ColorGray)

	// Set the key event handler to close the window on "Esc" or "q" key press
	statusBar.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Rune() == 'q' {
			app.Stop()
		}
		return event
	})

	addEntryBtn := controls.GetButton("Add Entry", func() {
		d.Pages.AddAndSwitchToPage(ADD_ENTRY, d.getAddEntyForm(), true)
	})

	addBulletBtn := controls.GetButton("Add Bullet", func() {
	})

	exitBtn := controls.GetButton("Exit", func() {
		app.Stop()
	})

	// Create a grid layout to hold the text view and status bar
	grid := tview.NewGrid().
		SetRows(-100, -5, -1).
		SetColumns(-1, -1, -1).
		AddItem(d.ContentTxtView, 0, 0, 1, 3, 1, 1, true).
		AddItem(addEntryBtn, 1, 0, 1, 1, 1, 1, false).
		AddItem(addBulletBtn, 1, 1, 1, 1, 1, 1, false).
		AddItem(exitBtn, 1, 2, 1, 1, 1, 1, false).
		AddItem(statusBar, 2, 0, 1, 3, 1, 1, false)

	tabItems := []interface{}{
		statusBar,
		addEntryBtn,
		addBulletBtn,
		exitBtn,
	}

	controls.NewTabControl(
		tabItems,
		grid,
		app,
	)

	d.Pages.AddPage("main", grid, true, true)

	app.SetRoot(d.Pages, true).SetFocus(statusBar)

	return app.Run()
}

func (d *TextModalWithQandEscLowerBar) mustGetTodaysFileText() string {
	content, err := directory.GetTodaysFileContent(d.FilePath)
	if err != nil {
		panic(fmt.Sprintf("Error reading file: %s\n", err))
	}
	return string(content)
}

func (d *TextModalWithQandEscLowerBar) getAddEntyForm() *tview.Form {
	value := ""
	form := tview.NewForm().
		AddInputField("Title", value, 40, nil, func(text string) {
			value = text
		}).
		AddTextView("Description", "This adds a header (#) to the markdown", 100, 1, true, false)

	form.SetBorder(true).SetTitle("Add Entry AKA dj add entry").SetTitleAlign(tview.AlignCenter)
	form.AddButton("Submit", func() {
		defer d.UpdateContent()

		if len(value) <= 0 {
			d.appModalToPage(ADD_ENTRY, "You gotta' have some text my dude")
			return
		}

		if len(strings.Trim(value, "\n")) == 0 {
			d.appModalToPage(ADD_ENTRY, "Blank lines are no good homie")
			value = ""
			form.GetFormItemByLabel("Title").(*tview.TextArea).SetText("", false)
			return
		}

		err := addlogic.AddEntryToFile(directory.GetTodaysFileName(d.FilePath), value)
		if err != nil {
			d.appModalToPage(ADD_ENTRY, err.Error())
		}

		d.Pages.RemovePage(ADD_ENTRY)
	})

	form.AddButton("Cancel", func() {
		d.Pages.RemovePage(ADD_ENTRY)
	})

	return form
}

func (d *TextModalWithQandEscLowerBar) appModalToPage(pageName string, msg string) {
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{"ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			d.Pages.RemovePage("modal")
		})

	d.Pages.AddPage("modal", modal, false, true)
}
