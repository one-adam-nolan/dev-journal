package display

import "io"

// Highlighter renders markdown content for terminal output.
type Highlighter interface {
	HighlightMarkdown(w io.Writer, content string) error
}
