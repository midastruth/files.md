package journal

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/Kunde21/markdownfmt/v3/markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"zakirullin/stuffbot/internal/fs"
)

var now = time.Now // to be replaced in tests

const (
	headerLevel        = 4
	intraNoteSeparator = "; "
)

func AddDailyNote(dir, noteFilename string, botFs *fs.FS, journalFilenameFormat, journalHeaderFormat string) error {
	content, err := botFs.Content(dir, noteFilename)
	if err != nil {
		return fmt.Errorf("failed to move to journal: can't get note content: %w", err)
	}
	note := fs.Title(noteFilename)
	if strings.TrimSpace(content) != "" {
		for _, line := range strings.Split(content, "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				note += intraNoteSeparator + line
			}
		}
	}
	journalFilename := now().Format(journalFilenameFormat)
	exists, err := botFs.Exists(fs.DirJournal, journalFilename)
	if err != nil {
		return err
	}
	if exists {
		content, err = botFs.Content(fs.DirJournal, journalFilename)
		if err != nil {
			return err
		}
	}
	content = insertDailyNote(content, journalHeaderFormat, note)
	return botFs.Put(fs.DirJournal, journalFilename, content)
}

func insertDailyNote(mdContent, journalHeaderFormat, note string) string {
	header := now().Format(journalHeaderFormat)
	r := markdown.NewRenderer()
	md := goldmark.New(
		goldmark.WithRenderer(r),
	)

	var buf bytes.Buffer

	source := []byte(mdContent)
	root := md.Parser().Parse(text.NewReader(source))
	root = addListItemAftreHeader(source, root, header, note)
	r.Render(&buf, source, root)
	return buf.String()
}

func addListItemAftreHeader(source []byte, root ast.Node, header, txt string) ast.Node {
	listItem := ast.NewListItem(0)
	listItem.AppendChild(listItem, ast.NewString([]byte(txt)))
	var nodeInserted bool

	ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		h, ok := node.(*ast.Heading)
		if !ok || !entering {
			return ast.WalkContinue, nil // skip all nodes except headings
		}
		headerText := h.Text(source)
		fmt.Println(string(headerText))
		if header != string(headerText) {
			return ast.WalkContinue, nil // it's not the header we are looking for
		}
		nodeInserted = true
		if list, ok := h.NextSibling().(*ast.List); ok {
			list.AppendChild(list, newListItem(txt))
		} else {
			h.InsertAfter(root, h, newList(newListItem(txt)))
		}
		return ast.WalkContinue, nil
	})
	if !nodeInserted {
		return appendNewSection(root, header, txt)
	}
	return root
}

func appendNewSection(root ast.Node, header, txt string) ast.Node {
	root.AppendChild(root, newHeader(header))
	root.AppendChild(root, newList(newListItem(txt)))
	return root
}

func newHeader(header string) *ast.Heading {
	heading := ast.NewHeading(headerLevel)
	heading.AppendChild(heading, ast.NewString([]byte(header)))
	return heading
}

func newList(listItem *ast.ListItem) *ast.List {
	list := ast.NewList('*')
	list.AppendChild(list, listItem)
	return list
}

func newListItem(txt string) *ast.ListItem {
	listItem := ast.NewListItem(0)
	listItem.AppendChild(listItem, ast.NewString([]byte(txt)))
	return listItem
}
