package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/app"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/middlewares"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/types"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// url для сокращения
type url struct {
	URL string `json:"url"`
}

// batchUrl в пакетной обработке
type batchUrl struct {
	CorrelationId string `json:"correlation_id"`
	OriginalUrl   string `json:"original_url"`
}

// shortenBatchUrl сокращенный урл в пакетной обработке
type shortenBatchUrl struct {
	CorrelationId string `json:"correlation_id"`
	ShortUrl      string `json:"short_url"`
}

// Сокращенный url
type response struct {
	URL string `json:"result"`
}

// URL пользователя
type userURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// CreateShortURLHandler — создает короткий урл.
func CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	h := md5.New()
	h.Write(body)

	hash := fmt.Sprintf("%x", h.Sum(nil))
	uuid := middlewares.UserSignedCookie.UUID
	shortURL := fmt.Sprintf("%s/%x", app.Cfg.BaseURL, h.Sum(nil))

	url := &types.URL{
		UUID:     uuid,
		Hash:     hash,
		URL:      string(body),
		ShortURL: shortURL,
	}

	err := app.Storage.Save(url)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(url.ShortURL))
}

// GetShortURLHandler — возвращает полный урл по короткому.
func GetShortURLHandler(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	exist, url, err := app.Storage.FindByHash(hash)

	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Location", url.URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte(url.URL))
	return
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
	uuid := middlewares.UserSignedCookie.UUID
	shortURL := fmt.Sprintf("%s/%x", app.Cfg.BaseURL, h.Sum(nil))

	url := &types.URL{
		UUID:     uuid,
		Hash:     hash,
		URL:      u.URL,
		ShortURL: shortURL,
	}

	err := app.Storage.Save(url)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	resp, _ := json.Marshal(response{URL: url.ShortURL})

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Accept", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// APICreateShortURLBatchHandler Api для создания коротких урлов пачками
func APICreateShortURLBatchHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData []batchUrl

	// Обрабатываем входящий json
	if err := json.NewDecoder(r.Body).Decode(&incomingData); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var urls []*types.URL
	var resp []*shortenBatchUrl
	uuid := middlewares.UserSignedCookie.UUID

	for _, url := range incomingData {
		shortURL := fmt.Sprintf("%s/%x", app.Cfg.BaseURL, url.CorrelationId)

		urls = append(urls, &types.URL{
			UUID:     uuid,
			Hash:     url.CorrelationId,
			URL:      url.OriginalUrl,
			ShortURL: shortURL,
		})
		resp = append(resp, &shortenBatchUrl{
			CorrelationId: url.CorrelationId,
			ShortUrl:      shortURL,
		})
	}

	err := app.Storage.SaveBatch(urls)
	if err != nil {
		fmt.Println(err)
	}

	response, _ := json.Marshal(resp)

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Accept", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

// GetUserURLSHandler — возвращает все сокращенные урлы пользователя.
func GetUserURLSHandler(w http.ResponseWriter, r *http.Request) {
	uuid := middlewares.UserSignedCookie.UUID

	urls, _ := app.Storage.FindByUUID(uuid)

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp := make([]userURL, 0, len(urls))

	for _, url := range urls {
		resp = append(resp, userURL{
			ShortURL:    url.ShortURL,
			OriginalURL: url.URL,
		})
	}

	respString, _ := json.Marshal(resp)

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Accept", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(respString)
}

// PingHandler проверяет соединение с базой
func PingHandler(w http.ResponseWriter, r *http.Request) {
	err := app.Storage.Ping()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("ok"))
}
