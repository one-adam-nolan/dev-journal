package tui

import (
	"dev-journal/directory"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type HistoryPage struct {
	Container  *tview.Flex
	FolderList *tview.List
	FileList   *tview.List
	Content    *tview.TextView
	Parent     Parent
	FilePath   string
}

func NewHistoryPage(parent Parent) *HistoryPage {
	return &HistoryPage{
		Parent:     parent,
		FilePath:   parent.GetFilePath(),
		FolderList: tview.NewList().ShowSecondaryText(false),
		FileList:   tview.NewList().ShowSecondaryText(false),
		Content:    tview.NewTextView(),
		Container:  tview.NewFlex(),
	}
}

func (hp *HistoryPage) Create() *tview.Flex {

	hp.setupFolderList()

	hp.setupFileList()

	hp.setupContent()

	hp.setupContainer()

	hp.addFoldersToList()

	return hp.Container
}

func (hp *HistoryPage) setupContainer() {
	hp.Container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			hp.Parent.GetPages().RemovePage(HISTORY)
		}

		return event
	})

	hp.Container.AddItem(hp.FolderList, 0, 1, true)
	hp.Container.AddItem(hp.FileList, 0, 1, false)
	hp.Container.AddItem(hp.Content, 0, 3, false)
}

func (hp *HistoryPage) setupContent() {
	hp.Content.SetBorder(true).SetTitle("Content")
	hp.Content.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			hp.Parent.GetApp().SetFocus(hp.FolderList)
		}
		return event
	})
}

func (hp *HistoryPage) setupFolderList() {
	hp.FolderList.SetBorder(true).SetTitle("Folders")
	hp.FolderList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		hp.FileList.Clear()

		currentFolder = mainText

		hp.refreshFileList()
	})

	hp.FolderList.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		hp.Parent.GetApp().SetFocus(hp.FileList)
	})
}

func (hp *HistoryPage) setupFileList() {
	hp.FileList.SetBorder(true).SetTitle("Files")
	hp.FileList.ShowSecondaryText(false)
	hp.FileList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		currentFile = mainText

		hp.updateContent()
	})

	hp.FileList.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		hp.Parent.GetApp().SetFocus(hp.Content)
	})
}

func (hp *HistoryPage) refreshFileList() {
	dirItems, err := directory.GetFolderContents(filepath.Join(hp.FilePath, currentFolder))
	if err != nil {
		log.Printf("%s", err)
		hp.FileList.AddItem(fmt.Sprintf("%s", err), "", 0, nil)
	}

	for _, d := range directory.SortDescendingByLastModified(dirItems) {
		if !d.IsDir() && !strings.HasSuffix(d.Name(), "swp") {
			hp.FileList.AddItem(d.Name(), d.Type().String(), 0, nil)
		}
	}
}

func (hp *HistoryPage) updateContent() {
	var fileContent string
	c, err := directory.GetContentForFile(filepath.Join(hp.FilePath, currentFolder, currentFile))
	if err != nil {
		fileContent = fmt.Sprintf("Error: %s", err)
	} else {
		fileContent = string(c)
	}

	hp.Content.SetText(fileContent)

	hp.Content.ScrollToBeginning()
}

func (hp *HistoryPage) addFoldersToList() {
	files, err := directory.GetFolderContents(hp.FilePath)
	if err != nil {
		fmt.Printf("%s \n", err)
	}

	for _, f := range directory.SortDescendingByLastModified(files) {
		if f.IsDir() && string(f.Name()[0]) != "." {
			hp.FolderList.AddItem(f.Name(), "", 0, func() {})
		}
	}
}
