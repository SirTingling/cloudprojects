package main

import (
	"log"
	"net/http"
	"os"

	"cloudprojects/calculator-backend-api/handlers"
	"github.com/rs/cors"
	"golang.org/x/exp/slog"
)

func main() {
	// Set up the logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Set up the HTTP server mux
	mux := http.NewServeMux()

	// Handlers to use with middleware
	mux.Handle("/add", handlers.LoggingMiddleware(http.HandlerFunc(handlers.AddHandler(logger)), logger))
	mux.Handle("/subtract", handlers.LoggingMiddleware(http.HandlerFunc(handlers.SubtractHandler(logger)), logger))
	mux.Handle("/multiply", handlers.LoggingMiddleware(http.HandlerFunc(handlers.MultiplyHandler(logger)), logger))
	mux.Handle("/divide", handlers.LoggingMiddleware(http.HandlerFunc(handlers.DivideHandler(logger)), logger))

	// Set up CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	// Wrap the mux with the CORS middleware
	handler := c.Handler(mux)

	// Start the server
	logger.Info("Starting server on :3000")
	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
