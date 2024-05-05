package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/MihailSergeenkov/shortener/internal/app/storage"
)

func AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h, err := storage.AddURL(string(body))

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	url := fmt.Sprintf(`http://localhost:8080/%s`, h)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(url))
}
