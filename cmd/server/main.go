package main

import (
	"cake-scraper/pkg/app"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := app.New(fiber.New())

	quit := make(chan os.Signal, 1)
	done := make(chan struct{}, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":3000"); err != nil {
			fmt.Println("Error starting server:", err)
			done <- struct{}{}
		}
	}()

	select {
	case <-quit:
		// Gracefully shutdown the server
		fmt.Println("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			fmt.Println("Error shutting down server:", err)
		}
		fmt.Println("Server exited.")
	case <-done:
		// Server exited
		fmt.Println("Server exited.")
	}
}
