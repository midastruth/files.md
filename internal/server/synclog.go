package server

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"zakirullin/stuffbot/config"
)

var lock sync.RWMutex

type LogEntry struct {
	Timestamp int64
	OldPath   string
	NewPath   string
}

func LogRename(oldPath, newPath string) {
	entry := LogEntry{
		Timestamp: time.Now().Unix(),
		OldPath:   oldPath,
		NewPath:   newPath,
	}

	lock.Lock()
	defer lock.Unlock()

	file, err := os.OpenFile(path.Join(config.BotCfg.WorkingDir, "fslog"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	record := fmt.Sprintf("%d\t%s\t%s\n", entry.Timestamp, entry.OldPath, entry.NewPath)

	_, err = file.WriteString(record)
	if err != nil {

	}
	err = file.Sync()
	if err != nil {

	}
}
