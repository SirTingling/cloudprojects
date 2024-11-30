package handlers

import (
	"encoding/json"
	"net/http"
	// "time"

	"golang.org/x/exp/slog"
)

type Request struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

type Response struct {
	Result float64 `json:"result"`
	Error  string  `json:"error,omitempty"`
}

func AddHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleOperation(logger, w, r, func(a, b float64) float64 { return a + b })
	}
}

func SubtractHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleOperation(logger, w, r, func(a, b float64) float64 { return a - b })
	}
}

func MultiplyHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleOperation(logger, w, r, func(a, b float64) float64 { return a * b })
	}
}

func DivideHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleOperation(logger, w, r, func(a, b float64) float64 {
			if b == 0 {
				panic("division by zero")
			}
			return a / b
		})
	}
}

func handleOperation(logger *slog.Logger, w http.ResponseWriter, r *http.Request, operation func(float64, float64) float64) {
	// Decode the request body
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		logger.Error("Invalid request payload", "error", err)
		return
	}

	// Perform the operation and send the response
	defer func() {
		if rec := recover(); rec != nil {
			logger.Error("Recovered from panic", "error", rec)
			http.Error(w, "Invalid operation: "+rec.(string), http.StatusBadRequest)
		}
	}()

	result := operation(req.A, req.B)
	resp := Response{Result: result}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		logger.Error("Failed to encode response", "error", err)
	}
}
