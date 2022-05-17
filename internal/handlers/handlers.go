package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Accept", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// CreateShortURLHandler — создает короткий урл.
func CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	sUrl := shortUrl(string(body))

	w.WriteHeader(http.StatusCreated)

	w.Write([]byte(sUrl))
}

// GetShortURLHandler — возвращает полный урл по короткому.
func GetShortURLHandler(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	url := getUrlByHash(hash)

	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)

	w.Write([]byte(url))
}

// shortUrl сокращает переданный url, сохраняет, возвращает короткую ссылку
func shortUrl(url string) (shortUrl string) {
	h := md5.New()
	h.Write([]byte(url))

	hash := fmt.Sprintf("%x", h.Sum(nil))
	shortUrl = fmt.Sprintf("%s/%x", os.Getenv("BASE_URL"), h.Sum(nil))

	urls[hash] = url
	return
}

// возвращает полный url по хешу
func getUrlByHash(hash string) (url string) {
	url = urls[hash]
	return
}
