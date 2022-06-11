package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/types"
	"io/ioutil"
	"net/http"

	"github.com/nastradamus39/ya_practicum_go_advanced/internal/app"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/middlewares"

	"github.com/go-chi/chi/v5"
)

// url для сокращения
type url struct {
	URL string `json:"url"`
}

// Сокращенный url
type response struct {
	URL string `json:"result"`
}

// Url пользователя
type userUrl struct {
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

// CreateShortURLHandler — создает короткий урл.
func CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	h := md5.New()
	h.Write(body)

	hash := fmt.Sprintf("%x", h.Sum(nil))
	uuid := middlewares.UserSignedCookie.Uuid
	shortUrl := fmt.Sprintf("%s/%x", app.Cfg.BaseURL, h.Sum(nil))

	url := &types.Url{
		Uuid:     uuid,
		Hash:     hash,
		URL:      string(body),
		ShortUrl: shortUrl,
	}

	err := app.Storage.Save(url)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(url.ShortUrl))
}

// GetShortURLHandler — возвращает полный урл по короткому.
func GetShortURLHandler(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	exist, url, err := app.Storage.FindByHash(hash)

	if !exist {
		w.WriteHeader(http.StatusNotFound)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("Location", url.URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte(url.URL))
}

// APICreateShortURLHandler Api для создания короткого урла
func APICreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	u := url{}

	// Обрабатываем входящий json
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h := md5.New()
	h.Write([]byte(u.URL))

	hash := fmt.Sprintf("%x", h.Sum(nil))
	uuid := middlewares.UserSignedCookie.Uuid
	shortUrl := fmt.Sprintf("%s/%x", app.Cfg.BaseURL, h.Sum(nil))

	url := &types.Url{
		Uuid:     uuid,
		Hash:     hash,
		URL:      u.URL,
		ShortUrl: shortUrl,
	}

	err := app.Storage.Save(url)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	resp, _ := json.Marshal(response{URL: url.ShortUrl})

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Accept", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// GetUserURLSHandler — возвращает все сокращенные урлы пользователя.
func GetUserURLSHandler(w http.ResponseWriter, r *http.Request) {
	uuid := middlewares.UserSignedCookie.Uuid

	urls, _ := app.Storage.FindByUuid(uuid)

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp := make([]userUrl, 0, len(urls))

	for _, url := range urls {
		resp = append(resp, userUrl{
			ShortUrl:    url.ShortUrl,
			OriginalUrl: url.URL,
		})
	}

	respString, _ := json.Marshal(resp)

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Accept", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(respString)
}
