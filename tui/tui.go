package tui

type Displayable interface {
	Display() error
}

func DisplayTodayModal(content string) Displayable {
	return &TextModalWithQandEscLowerBar{
		Content: content,
	}
}
