package logger

import (
	"cake-scraper/pkg/util"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

var logPath = filepath.Join(util.ProjectRoot, "log/scraper.log")

func init() {
	_ = os.MkdirAll(filepath.Dir(logPath), 0755)
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	slog.SetDefault(slog.New(
		slog.NewJSONHandler(file, &slog.HandlerOptions{
			AddSource: true,
		}),
	))
}
