package i18n

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func LoadEmojiFile(path string) error {
	emojiFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("i18n.LoadEmojiFile: %w", err)
	}
	defer emojiFile.Close()

	bytes, err := io.ReadAll(emojiFile)
	if err != nil {
		return fmt.Errorf("i18n.LoadEmojiFile: %w", err)
	}

	var emojis map[string][]string
	err = json.Unmarshal(bytes, &emojis)
	if err != nil {
		return fmt.Errorf("i18n.LoadEmojiFile: can't unmarshal: %w", err)
	}

	emojisByKeyword = make(map[string]string)
	for emoji, keywords := range emojis {
		for _, keyword := range keywords {
			emojisByKeyword[keyword] = emoji
		}
	}

	return nil
}

func Emojify(str string) string {
	strLower := strings.ToLower(str)
	aliases := []string{strLower, strLower + "s", strings.TrimSuffix(strLower, "s")}
	for _, alias := range aliases {
		icon, _ := emojisByKeyword[alias]
		if icon != "" {
			return fmt.Sprintf("%s %s", icon, str)
		}
	}

	for _, word := range strings.Fields(str) {
		icon, _ := emojisByKeyword[word]
		if icon != "" {
			return fmt.Sprintf("%s %s", icon, str)
		}
	}

	return str
}
