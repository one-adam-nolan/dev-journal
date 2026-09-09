package journal

import "errors"

var (
	ErrJournalNotFound  = errors.New("journal entry file not found")
	ErrDayAlreadyExists = errors.New("journal entry already exists for today")
	ErrInvalidDate      = errors.New("invalid date format, use MM/DD/YYYY")
)
