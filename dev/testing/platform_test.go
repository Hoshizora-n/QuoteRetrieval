package testing

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func QuoteHandler(w http.ResponseWriter, r *http.Request) {
	var response interface{}
	response = []map[string]interface{}{
		{
			"quote":    "if you don't fight, you won't win",
			"author":   "eren yeager",
			"category": "freedom",
		},
	}
	status := http.StatusOK

	apiKey := r.Header.Get("X-API-KEY")
	if apiKey == "" {
		response = map[string]interface{}{
			"error": "missing api key",
		}
	}

	jsonResponse, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonResponse)
}

func PlatformServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		QuoteHandler(w, r)
	}))
	return server
}
