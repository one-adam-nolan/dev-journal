package tui

import (
	"dev-journal/directory"
	"dev-journal/pkg/addlogic"

	"github.com/rivo/tview"
)

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
			d.appModalToPage("You gotta' have some text my dude")
			return
		}

		err := addlogic.AddEntryToFile(directory.GetTodaysFileName(d.FilePath), value)
		if err != nil {
			d.appModalToPage(err.Error())
		}

		d.Pages.RemovePage(ADD_ENTRY)
	})

	form.AddButton("Cancel", func() {
		d.Pages.RemovePage(ADD_ENTRY)
	})

	return form
}
