package i18n

const (
	emojisPath = "assets/emoji.json"
)

func Emojify(str string) string {
	// TODO load emojis from file
	if str == "Took a break" {
		return "🍅 Took a break"
	}

	if str == "Today" {
		return "🏠 Today"
	}

	return str
}
