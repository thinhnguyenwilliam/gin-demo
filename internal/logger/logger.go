// gin-demo/internal/logger/logger.go
package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var (
	Log  zerolog.Logger
	file *os.File
	mu   sync.Mutex
)

func Init() {
	_ = os.MkdirAll("logs", os.ModePerm)

	createNewLogFile()

	go scheduleMidnightRotation()
	go cleanupOldLogs()
}

func createNewLogFile() {
	mu.Lock()
	defer mu.Unlock()

	if file != nil {
		_ = file.Close()
	}

	filename := fmt.Sprintf("logs/app-%s.log",
		time.Now().Format("2006-01-02"))

	f, err := os.OpenFile(
		filename,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		panic(err)
	}

	file = f

	multi := io.MultiWriter(os.Stdout, file)

	Log = zerolog.New(multi).
		With().
		Timestamp().
		Logger()
}

func scheduleMidnightRotation() {
	for {
		now := time.Now()

		// Next midnight (00:00)
		next := time.Date(
			now.Year(),
			now.Month(),
			now.Day()+1,
			0, 0, 0, 0,
			now.Location(),
		)

		time.Sleep(time.Until(next))

		createNewLogFile()
	}
}

func cleanupOldLogs() {
	for {
		time.Sleep(24 * time.Hour)

		files, _ := filepath.Glob("logs/*.log")
		cutoff := time.Now().AddDate(0, 0, -90)

		for _, f := range files {
			info, err := os.Stat(f)
			if err != nil {
				continue
			}

			if info.ModTime().Before(cutoff) {
				_ = os.Remove(f)
			}
		}
	}
}
