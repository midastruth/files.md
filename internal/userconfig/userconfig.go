// Package userconfig stores user's configuration in file.
// It stores such settings for users as: language, home, quick buttons, schedule and so on.
package userconfig

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/exp/slog"

	"zakirullin/dumpbot/i18n"
)

var DefaultConfig = Config{
	config: config{
		Language:               "en",
		HomeCmd:                "today",
		MoveToButtons:          []string{"tomorrow", "later", "day", "note", "checklist", "doc", "recent", "journal"},
		PomodoroDurationMinute: 25,
	},
}

var TasksOnlyConfig = Config{
	config: config{
		HomeCmd:       "today",
		MoveToButtons: []string{"tomorrow", "later", "day"},
	},
}

var NotesOnlyConfig = Config{
	config: config{
		HomeCmd:       "notes",
		MoveToButtons: []string{"##NOTE_DIRS##"},
	},
}

type Config struct {
	config
}

type config struct {
	Language               string   `json:"language"`
	HomeCmd                string   `json:"homeCmd"`
	MoveToButtons          []string `json:"moveToButtons"`
	PomodoroDurationMinute float64  `json:"pomodoroDurationMinute"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &c.config)
}

// TODO add file creation
func (c *Config) LoadOrCreate(path string) error {
	configFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("config.LoadOrCreate: %w", err)
	}
	defer configFile.Close()

	bytes, err := io.ReadAll(configFile)
	if err != nil {
		return fmt.Errorf("config.LoadOrCreate: %w", err)
	}

	err = json.Unmarshal(bytes, c)
	if err != nil {
		return fmt.Errorf("config.LoadOrCreate: can't unmarshal: %w", err)
	}

	return nil
}

func (c *Config) Save(path string) {

}

func mapConfigButtonNamesToRealNames(configNames []string) []string {
	configToReal := map[string]string{
		"tomorrow":  i18n.StrForTomorrow,
		"later":     i18n.StrForLater,
		"day":       i18n.StrForDay,
		"note":      i18n.StrToNote,
		"checklist": i18n.StrToChecklist,
		"doc":       i18n.StrToDoc,
	}

	var realNames []string
	for _, configName := range configNames {
		realName, ok := configToReal[configName]
		if !ok {
			continue
		}

		realNames = append(realNames, realName)
	}

	return realNames
}

func (c *Config) SetPomodoroDuration(value time.Duration) error {
	if value <= 0 || value > 24*time.Hour {
		return fmt.Errorf("config.SetPomodoroDuration: value is invalid: %v", value)
	}
	c.config.PomodoroDurationMinute = value.Minutes()
	return nil
}

func (c *Config) PomodoroDuration() time.Duration {
	minutes := c.config.PomodoroDurationMinute
	if minutes <= 0 {
		slog.Error("Pomodoro duration is invalid. Using default value", "duration",
			c.config.PomodoroDurationMinute, "default", DefaultConfig.config.PomodoroDurationMinute)
		//I don't use DefaultConfig.PomodoroDuration() because it may cause infinite recursion
		minutes = DefaultConfig.config.PomodoroDurationMinute
	}
	return time.Duration(minutes * float64(time.Minute))
}
