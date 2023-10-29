package tui

import (
	"dev-journal/directory"
	"dev-journal/pkg/controls"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	MAIN       = "main"
	ADD_ENTRY  = "input-entry"
	ADD_BULLET = "input-bulet"
	HISTORY    = "history"
	MODAL      = "modal"
)

var (
	currentFolder string = ""
	currentFile   string = ""
)

type Parent interface {
	GetFilePath() string
	GetPages() *tview.Pages
	GetContentView() *tview.TextView
	GetApp() *tview.Application
}

type TextModalWithQandEscLowerBar struct {
	FilePath       string
	Pages          *tview.Pages
	ContentTxtView *tview.TextView
	App            *tview.Application
}

func (d *TextModalWithQandEscLowerBar) UpdateContent() {
	d.ContentTxtView.SetText(d.mustGetTodaysFileText())
}

func (d *TextModalWithQandEscLowerBar) GetFilePath() string {
	return d.FilePath
}

func (d *TextModalWithQandEscLowerBar) GetPages() *tview.Pages {
	return d.Pages
}

func (d *TextModalWithQandEscLowerBar) GetContentView() *tview.TextView {
	return d.ContentTxtView
}

func (d *TextModalWithQandEscLowerBar) GetApp() *tview.Application {
	return d.App
}

func (d *TextModalWithQandEscLowerBar) Display() error {
	d.App = tview.NewApplication()
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
		SetDynamicColors(true).
		SetTextStyle(tcell.StyleDefault.
			Foreground(tcell.Color111).
			Background(tcell.Color234))

	d.UpdateContent()

	addEntryBtn := d.getNavigationButton("Add Entry", ADD_ENTRY, d.getAddEntyForm)

	addBulletBtn := d.getNavigationButton("Add Bullet", ADD_BULLET, d.getAddBulletForm)

	exitBtn := controls.GetButton("Exit", func() {
		d.App.Stop()
	})

	historyBtn := controls.GetButton("History", func() {
		d.Pages.AddAndSwitchToPage(HISTORY, NewHistoryPage(d).Create(), true)
	})

	// Create a grid layout to hold the text view and status bar
	grid := tview.NewGrid().
		SetRows(-100, 3, 3).
		SetColumns(-1, -1, -1).
		AddItem(d.ContentTxtView, 0, 0, 1, 3, 1, 1, true).
		AddItem(addEntryBtn, 1, 0, 1, 1, 10, 1, false).
		AddItem(addBulletBtn, 1, 1, 1, 1, 10, 1, false).
		AddItem(exitBtn, 1, 2, 1, 1, 1, 10, false).
		AddItem(historyBtn, 2, 0, 1, 3, 1, 10, false)

	tabItems := []interface{}{
		d.ContentTxtView,
		addEntryBtn,
		addBulletBtn,
		exitBtn,
		historyBtn,
	}

	controls.NewTabControl(
		tabItems,
		grid,
		d.App,
	)

	d.Pages.AddPage(MAIN, grid, true, true)

	d.App.SetRoot(d.Pages, true).SetFocus(d.ContentTxtView)

	return d.App.Run()
}

func (d *TextModalWithQandEscLowerBar) getNavigationButton(title string, tag string, targetFunc func() *tview.Form) *tview.Button {
	return controls.GetButton(title, func() {
		d.Pages.AddAndSwitchToPage(tag, targetFunc(), true)
	})
}

func (d *TextModalWithQandEscLowerBar) mustGetTodaysFileText() string {
	content, err := directory.GetTodaysFileContent(d.FilePath)
	if err != nil {
		panic(fmt.Sprintf("Error reading file: %s\n", err))
	}
	return string(content)
}

func (d *TextModalWithQandEscLowerBar) appModalToPage(msg string) {
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{"ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			d.Pages.RemovePage(MODAL)
		})

	d.Pages.AddPage(MODAL, modal, false, true)
}
