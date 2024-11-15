package logger

import (
	"bufio"
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
	w := bufio.NewWriter(file)
	slog.SetDefault(slog.New(
		slog.NewJSONHandler(w, &slog.HandlerOptions{
			AddSource: true,
		}),
	))
}
