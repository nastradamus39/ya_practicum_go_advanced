package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var urls = map[string]string{}

// url для сокращения
type url struct {
	Url string `json:"url"`
}

// Сокращенный url
type response struct {
	Url string `json:"result"`
}

// ApiCreateShortURLHandler создает короткий урл
func ApiCreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	url := url{}

	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url.Url = shortUrl(url.Url)

	resp, _ := json.Marshal(response(url))

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

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

// shortUrl сокращает переданный url
func shortUrl(url string) (shortUrl string) {
	h := md5.New()
	return fmt.Sprintf("http://127.0.0.1:8080/%x", h.Sum(nil))
}
