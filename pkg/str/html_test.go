package str

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarkdownToHtmlHeader(t *testing.T) {
	r := require.New(t)

	md := `# Header`
	html := MarkdownToHtml(md)

	r.Equal("<b>Header</b>", html)
}

func TestMarkdownToHtmlHeaderAndText(t *testing.T) {
	r := require.New(t)

	md := "# Header\nText"
	html := MarkdownToHtml(md)

	r.Equal("<b>Header</b>\n\nText", html)
}

func TestMarkdownToHtmlBold(t *testing.T) {
	r := require.New(t)

	md := `**Bold**`
	html := MarkdownToHtml(md)

	r.Equal("<strong>Bold</strong>", html)
}
