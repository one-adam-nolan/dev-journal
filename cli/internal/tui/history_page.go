package tui

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type HistoryPage struct {
	Container  *tview.Flex
	FolderList *tview.List
	FileList   *tview.List
	Content    *tview.TextView
	Parent     Parent
}

func NewHistoryPage(parent Parent) *HistoryPage {
	return &HistoryPage{
		Parent:     parent,
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
	ctx := context.Background()
	resp, err := hp.Parent.GetJournal().ListFiles(ctx, currentFolder)
	if err != nil {
		hp.FileList.AddItem(fmt.Sprintf("%s", err), "", 0, nil)
		return
	}

	for _, name := range resp.GetFileNames() {
		hp.FileList.AddItem(name, "", 0, nil)
	}
}

func (hp *HistoryPage) updateContent() {
	ctx := context.Background()
	resp, err := hp.Parent.GetJournal().ReadFile(ctx, currentFolder, currentFile)
	if err != nil {
		hp.Content.SetText(fmt.Sprintf("Error: %s", err))
	} else {
		hp.Content.SetText(string(resp.GetContent()))
	}

	hp.Content.ScrollToBeginning()
}

func (hp *HistoryPage) addFoldersToList() {
	ctx := context.Background()
	resp, err := hp.Parent.GetJournal().ListMonthFolders(ctx)
	if err != nil {
		fmt.Printf("%s \n", err)
		return
	}

	for _, name := range resp.GetFolderNames() {
		hp.FolderList.AddItem(name, "", 0, func() {})
	}
}
