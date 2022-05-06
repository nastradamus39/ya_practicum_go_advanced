package handlers

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var urls = map[string]string{}

// CreateShortURLHandler — создает короткий урл.
func CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	url := string(body)

	defer r.Body.Close()

	h := md5.New()
	h.Write(body)

	hash := fmt.Sprintf("%x", h.Sum(nil))

	urls[hash] = url

	w.WriteHeader(http.StatusCreated)

	w.Write([]byte(fmt.Sprintf("http://127.0.0.1:8080/%s", hash)))
}

// GetShortURLHandler — возвращает полный урл по короткому.
func GetShortURLHandler(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	url := urls[hash]

	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)

	w.Write([]byte(url))
}
