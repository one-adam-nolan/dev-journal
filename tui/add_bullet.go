package tui

import (
	"dev-journal/directory"
	"dev-journal/pkg/addlogic"
	"strings"

	"github.com/rivo/tview"
)

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
			d.appModalToPage("You gotta' have some text my dude")
			return
		}

		value = strings.Trim(value, "\n")

		if len(value) == 0 {
			d.appModalToPage("Blank lines are no good homie")
			value = ""
			form.GetFormItemByLabel("Message").(*tview.TextArea).SetText("", false)
			return
		}

		err := addlogic.AddBulletToFile(directory.GetTodaysFileName(d.FilePath), value)
		if err != nil {
			d.appModalToPage(err.Error())
		}

		d.Pages.RemovePage(ADD_BULLET)
	})

	form.AddButton("Cancel", func() {
		d.Pages.RemovePage(ADD_BULLET)
	})

	return form
}
