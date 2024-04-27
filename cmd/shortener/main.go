package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

const ADDR = "localhost:8080"

var uMap = make(map[string]string, 100)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(ADDR, http.HandlerFunc(webhook))
}

func webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost && r.URL.Path == "/" {
		processPost(w, r)
		return
	}

	if r.Method == http.MethodGet {
		processGet(w, r)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func processPost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h, err := randomHex()

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	uMap[h] = string(body)
	url := fmt.Sprintf(`http://%s/%s`, ADDR, h)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(url))
}

func randomHex() (string, error) {
	bytes := make([]byte, 4)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func processGet(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile(`^/\w{8}$`)
	key := re.FindString(r.URL.Path)

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, present := uMap[strings.TrimLeft(key, "/")]

	if !present {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", u)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
