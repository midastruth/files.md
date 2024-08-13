package txt

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

	r.Equal("<b>Header</b>\nText", html)
}

func TestMarkdownToHtmlBold(t *testing.T) {
	r := require.New(t)

	md := "**bold**"
	html := MarkdownToHtml(md)

	r.Equal("<b>bold</b>", html)
}

func TestMarkdownToHtmlItalic(t *testing.T) {
	r := require.New(t)

	md := "*italic*"
	html := MarkdownToHtml(md)

	r.Equal("<i>italic</i>", html)
}

func TestMarkdownToHtmlInvalid(t *testing.T) {
	r := require.New(t)

	md := "__valid__**invalid"
	html := MarkdownToHtml(md)

	r.Equal("<b>valid</b>**invalid", html)
}

func TestMarkdownToHtmlMultiline(t *testing.T) {
	r := require.New(t)

	md := "line1\n**line2**\nline3"
	html := MarkdownToHtml(md)

	r.Equal("line1\n<b>line2</b>\nline3", html)
}

func TestMarkdownToHtmlNoLists(t *testing.T) {
	r := require.New(t)

	md := "list\n1) item1\n2) item2"
	html := MarkdownToHtml(md)

	r.Equal("list\n1) item1\n2) item2", html)
}
