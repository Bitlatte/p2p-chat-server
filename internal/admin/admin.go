package admin

import (
	"net/http"
	"os"
)

var adminAPIKey = os.Getenv("ADMIN_API_KEY")

func IsAuthenticated(r *http.Request) bool {
	apiKey := r.Header.Get("X-API-Key")
	return apiKey == adminAPIKey
}
