package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"net/http"
)

var BaseUrl string

var Storage *FileStorage

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

	u, err := getUrlByHash(hash)

	if err != nil {
		fmt.Println(fmt.Sprintf("Cannot find full url. Error - %s", err))
	}

	w.Header().Add("Location", u)
	w.WriteHeader(http.StatusTemporaryRedirect)

	w.Write([]byte(u))
}

// shortUrl сокращает переданный url, сохраняет, возвращает короткую ссылку
func shortUrl(url string) (shortUrl string) {
	h := md5.New()
	h.Write([]byte(url))

	hash := fmt.Sprintf("%x", h.Sum(nil))

	u, _ := Storage.Find(hash)
	if u == "" {
		// Сохраняем на диск
		Storage.Save(hash, url)
	}

	urls[hash] = url // сохраняем в памяти

	shortUrl = fmt.Sprintf("%s/%x", BaseUrl, h.Sum(nil))

	return
}

// возвращает полный url по хешу
func getUrlByHash(hash string) (url string, err error) {
	// Ищем в памяти
	u := urls[hash]
	if u != "" {
		return u, nil
	}

	// Если в памяти нет - ищем в файле
	if u == "" {
		u, err = Storage.Find(hash)
	}
	return url, err
}
