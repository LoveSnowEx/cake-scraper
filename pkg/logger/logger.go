package logger

import (
	"log"
	"log/slog"
	"os"
)

func init() {
	_ = os.MkdirAll("log", 0755)
	file, err := os.OpenFile("log/scraper.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	slog.SetDefault(slog.New(
		slog.NewJSONHandler(file, &slog.HandlerOptions{
			AddSource: true,
		}),
	))
}
