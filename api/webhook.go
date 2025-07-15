package api

import (
	"net/http"
	"tg-remote/internal/api"
)

// Handler is the main webhook handler for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Delegate to the internal webhook handler
	api.Handler(w, r)
}
