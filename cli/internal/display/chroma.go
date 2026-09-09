package display

import (
	"io"

	"github.com/alecthomas/chroma/quick"
)

// ChromaHighlighter highlights markdown using chroma.
type ChromaHighlighter struct{}

// NewChromaHighlighter returns a chroma-backed markdown highlighter.
func NewChromaHighlighter() *ChromaHighlighter {
	return &ChromaHighlighter{}
}

func (h *ChromaHighlighter) HighlightMarkdown(w io.Writer, content string) error {
	return quick.Highlight(w, content, "markdown", "terminal16m", "monokai")
}
