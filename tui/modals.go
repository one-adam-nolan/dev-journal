package tui

import (
	"dev-journal/directory"
	"dev-journal/pkg/addlogic"
	"dev-journal/pkg/controls"
	"fmt"
	"path/filepath"
	"strings"

	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	MAIN       = "main"
	ADD_ENTRY  = "input-entry"
	ADD_BULLET = "input-bulet"
	HISTORY    = "history"
)

var (
	currentFolder string = ""
	currentFile   string = ""
)

type TextModalWithQandEscLowerBar struct {
	FilePath       string
	Pages          *tview.Pages
	ContentTxtView *tview.TextView
	App            *tview.Application
}

func (d *TextModalWithQandEscLowerBar) UpdateContent() {
	d.ContentTxtView.SetText(d.mustGetTodaysFileText())
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

	addEntryBtn := controls.GetButton("Add Entry", func() {
		d.Pages.AddAndSwitchToPage(ADD_ENTRY, d.getAddEntyForm(), true)
	})

	addBulletBtn := controls.GetButton("Add Bullet", func() {
		d.Pages.AddAndSwitchToPage(ADD_BULLET, d.getAddBulletForm(), true)
	})

	exitBtn := controls.GetButton("Exit", func() {
		d.App.Stop()
	})

	historyBtn := controls.GetButton("History", func() {
		d.Pages.AddAndSwitchToPage(HISTORY, d.getHistoryPage(), true)
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

func (d *TextModalWithQandEscLowerBar) getAddBulletForm() *tview.Form {
	value := ""
	form := tview.NewForm().
		AddTextArea("Message", value, 100, 10, 500, func(text string) {
			value = text
		}).
		AddTextView("Description", "This adds a bullet underneath a header", 100, 1, true, false)

	form.SetBorder(true).SetTitle("Add Entry AKA dj add entry").SetTitleAlign(tview.AlignCenter)
	form.AddButton("Submit", func() {
		defer d.UpdateContent()

		if len(value) <= 0 {
			d.appModalToPage(ADD_BULLET, "You gotta' have some text my dude")
			return
		}

		if len(strings.Trim(value, "\n")) == 0 {
			d.appModalToPage(ADD_BULLET, "Blank lines are no good homie")
			value = ""
			form.GetFormItemByLabel("Message").(*tview.TextArea).SetText("", false)
			return
		}

		err := addlogic.AddBulletToFile(directory.GetTodaysFileName(d.FilePath), value)
		if err != nil {
			d.appModalToPage(ADD_BULLET, err.Error())
		}

		d.Pages.RemovePage(ADD_BULLET)
	})

	form.AddButton("Cancel", func() {
		d.Pages.RemovePage(ADD_BULLET)
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

func (d *TextModalWithQandEscLowerBar) getHistoryPage() *tview.Flex {
	flex := tview.NewFlex()

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			d.Pages.RemovePage(HISTORY)
		}

		return event
	})

	folderList := tview.NewList().ShowSecondaryText(false)
	folderList.SetBorder(true).SetTitle("Folders")

	files := tview.NewList().ShowSecondaryText(false)
	files.SetBorder(true).SetTitle("Files")
	files.ShowSecondaryText(false)

	content := tview.NewTextView()
	content.SetBorder(true).SetTitle("Content")

	folderList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		files.Clear()

		currentFolder = mainText

		dirItems, err := directory.GetFolderContents(filepath.Join(d.FilePath, currentFolder))
		if err != nil {
			log.Printf("%s", err)
			files.AddItem(fmt.Sprintf("%s", err), "", 0, nil)
		}

		for _, d := range directory.SortDescendingByLastModified(dirItems) {
			if !d.IsDir() && !strings.HasSuffix(d.Name(), "swp") {
				files.AddItem(d.Name(), d.Type().String(), 0, nil)
			}
		}
	})

	folderList.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		d.App.SetFocus(files)
	})

	files.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		currentFile = mainText

		var fileContent string
		c, err := directory.GetContentForFile(filepath.Join(d.FilePath, currentFolder, currentFile))
		if err != nil {
			fileContent = fmt.Sprintf("Error: %s", err)
		} else {
			fileContent = string(c)
		}

		content.SetText(fileContent)
	})

	files.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		d.App.SetFocus(content)
	})

	content.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			d.App.SetFocus(folderList)
		}
		return event
	})

	flex.AddItem(folderList, 0, 1, true)
	flex.AddItem(files, 0, 1, false)
	flex.AddItem(content, 0, 3, false)

	item, err := directory.GetFolderContents(d.FilePath)
	if err != nil {
		fmt.Printf("%s \n", err)
	}

	for _, f := range directory.SortDescendingByLastModified(item) {
		if f.IsDir() && string(f.Name()[0]) != "." {
			folderList.AddItem(f.Name(), "", 0, func() {})
		}
	}

	return flex
}
