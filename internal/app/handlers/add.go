package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/MihailSergeenkov/shortener/internal/app/storage"
)

func AddHandler(w http.ResponseWriter, r *http.Request) {
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

	url := fmt.Sprintf(`%s/%s`, config.Params.UAddr, h)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(url))
}
