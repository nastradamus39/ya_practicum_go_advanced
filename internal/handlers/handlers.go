package handlers

import (
	"encoding/json"
	"fmt"
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

	url := app.ShortURL(string(body))

	w.WriteHeader(http.StatusCreated)

	w.Write([]byte(url.ShortUrl))
}

// GetShortURLHandler — возвращает полный урл по короткому.
func GetShortURLHandler(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	uuid := middlewares.UserSignedCookie.Uuid

	u, err := app.GetURLByHash(uuid, hash)

	if err != nil {
		fmt.Printf("Cannot find full url. Error - %s", err)
	}

	w.Header().Add("Location", u.URL)
	w.WriteHeader(http.StatusTemporaryRedirect)

	w.Write([]byte(u.URL))
}

// APICreateShortURLHandler Api для создания короткого урла
func APICreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	url := url{}

	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url.URL = app.ShortURL(url.URL).ShortUrl

	resp, _ := json.Marshal(response(url))

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Accept", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// GetUserURLSHandler — возвращает все сокращенные урлы пользователя.
func GetUserURLSHandler(w http.ResponseWriter, r *http.Request) {
	uuid := middlewares.UserSignedCookie.Uuid

	urls := app.GetUrlsByUuid(uuid)
	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(""))
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
