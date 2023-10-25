package tui

type Displayable interface {
	Display() error
}

func DisplayTodayModal(filePath string) Displayable {
	return &TextModalWithQandEscLowerBar{
		FilePath: filePath,
	}
}
