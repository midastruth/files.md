package txt

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/require"
)

func TestBold(t *testing.T) {
	r := require.New(t)

	text := "bold"
	messageEntities := []tgbotapi.MessageEntity{
		{Type: "bold", Offset: 0, Length: 4},
	}

	md := TelegramEntitiesToMarkdown(text, messageEntities)
	r.Equal("**bold**", md)
}

func TestItalic(t *testing.T) {
	r := require.New(t)

	text := "italic"
	messageEntities := []tgbotapi.MessageEntity{
		{Type: "italic", Offset: 0, Length: 6},
	}

	md := TelegramEntitiesToMarkdown(text, messageEntities)
	r.Equal("*italic*", md)
}

func TestBoldAndItalic(t *testing.T) {
	r := require.New(t)

	text := "BoldAndItalic"
	messageEntities := []tgbotapi.MessageEntity{
		{Type: "bold", Offset: 0, Length: 13},
		{Type: "italic", Offset: 0, Length: 13},
	}

	md := TelegramEntitiesToMarkdown(text, messageEntities)
	r.Equal("***BoldAndItalic***", md)
}

func TestBoldThenItalic(t *testing.T) {
	r := require.New(t)

	text := "bolditalic"
	messageEntities := []tgbotapi.MessageEntity{
		{Type: "bold", Offset: 0, Length: 4},
		{Type: "italic", Offset: 4, Length: 6},
	}

	md := TelegramEntitiesToMarkdown(text, messageEntities)
	r.Equal("**bold***italic*", md)
}

func TestLink(t *testing.T) {
	r := require.New(t)

	text := "l"
	messageEntities := []tgbotapi.MessageEntity{
		{Type: "text_link", Offset: 0, Length: 1, URL: "google.com"},
	}

	md := TelegramEntitiesToMarkdown(text, messageEntities)
	r.Equal("[l](google.com)", md)
}

func TestMultilineTextWithMarkdown(t *testing.T) {
	r := require.New(t)

	text := "header\nitalic\n\nAlso italic\n\nheader2\nitalic\ncode"
	messageEntities := []tgbotapi.MessageEntity{
		{Type: "bold", Offset: 0, Length: 7},
		{Type: "italic", Offset: 7, Length: 21},
		{Type: "bold", Offset: 28, Length: 8},
		{Type: "italic", Offset: 36, Length: 7},
		{Type: "code", Offset: 43, Length: 4},
	}

	markdown := TelegramEntitiesToMarkdown(text, messageEntities)
	expectedMarkdown := "**header**\n*italic*\n\n*Also italic*\n\n**header2**\n*italic*\n`code`"
	r.Equal(expectedMarkdown, markdown)
}

func TestSpacedItalic(t *testing.T) {
	r := require.New(t)
	text := "Header\nLeverage one Minute Praising instead"

	messageEntities := []tgbotapi.MessageEntity{
		{Type: "italic", Offset: 16, Length: 20},
	}

	markdown := TelegramEntitiesToMarkdown(text, messageEntities)
	expectedMarkdown := "Header\nLeverage *one Minute Praising* instead"
	r.Equal(expectedMarkdown, markdown)
}

func TestEmojiInMessageEntities(t *testing.T) {
	r := require.New(t)

	text := "👍b"
	messageEntities := []tgbotapi.MessageEntity{
		{Type: "bold", Offset: 2, Length: 1}, // Emoji is 4 bytes or 2 runes
	}

	md := TelegramEntitiesToMarkdown(text, messageEntities)
	r.Equal("👍**b**", md)
}

func TestSkinEmoji(t *testing.T) {
	r := require.New(t)

	text := "🤘🏾b"
	messageEntities := []tgbotapi.MessageEntity{
		{Type: "bold", Offset: 4, Length: 1}, // Tone emoji is 8 bytes or 4 runes
	}

	md := TelegramEntitiesToMarkdown(text, messageEntities)
	r.Equal("🤘🏾**b**", md)
}

func TestPre(t *testing.T) {
	r := require.New(t)

	text := "line1\nline2"
	messageEntities := []tgbotapi.MessageEntity{
		{Type: "pre", Offset: 0, Length: 11},
	}

	md := TelegramEntitiesToMarkdown(text, messageEntities)
	r.Equal("```line1\nline2```", md)
}

func TestDoesntEscapeMD(t *testing.T) {
	r := require.New(t)

	text := "Ask @_a_ __b__ *a* **b** `c` ```multiline```"
	md := TelegramEntitiesToMarkdown(text, nil)
	r.Equal("Ask @_a_ __b__ *a* **b** `c` ```multiline```", md)
}

func TestDoesntEscapeBrokenMD(t *testing.T) {
	r := require.New(t)

	text := "Ask @nick_name * `"
	md := TelegramEntitiesToMarkdown(text, nil)
	r.Equal("Ask @nick_name * `", md)

	text = "___ *** __ ```"
	md = TelegramEntitiesToMarkdown(text, nil)
	r.Equal("___ *** __ ```", md)
}
