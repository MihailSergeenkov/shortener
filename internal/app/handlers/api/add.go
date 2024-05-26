package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/models"
)

func APIAddHandler(s data.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.Request
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to read request body: %v", err)
			return
		}

		u, err := s.AddURL(req.URL)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to add URL to storage: %v", err)
			return
		}

		result, err := url.JoinPath(config.Params.BaseURL, u.ShortUrl)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to construct URL: %v", err)
			return
		}

		resp := models.Response{Result: result}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error encoding responseL: %v", err)
			return
		}
	}
}
