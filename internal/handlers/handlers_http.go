package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nastradamus39/ya_practicum_go_advanced/internal/app"
	shortenerErrors "github.com/nastradamus39/ya_practicum_go_advanced/internal/errors"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/middlewares"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/types"

	"github.com/go-chi/chi/v5"
)

// url для сокращения
type url struct {
	URL string `json:"url"`
}

// batchURL в пакетной обработке
type batchURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// shortenBatchURL сокращенный урл в пакетной обработке
type shortenBatchURL struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
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

// CreateShortURLHTTPHandler — создает короткий урл.
func CreateShortURLHTTPHandler(w http.ResponseWriter, r *http.Request) {
	originalURL, _ := ioutil.ReadAll(r.Body)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("CreateShortURLHandler. %s", err)
		}
	}(r.Body)

	uuid := middlewares.UserSignedCookie.UUID
	url, err := CreateShortURLHandler(string(originalURL), uuid)

	// Если такой url уже есть - отдаем соответствующий статус
	if errors.Is(err, shortenerErrors.ErrURLConflict) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(url.ShortURL))
		return
	}

	// Другие ошибки при сохранении в хранилище
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(url.ShortURL))
}

// GetShortURLHTTPHandler — возвращает полный урл по короткому.
func GetShortURLHTTPHandler(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	url, err := GetShortURLHandler(hash)

	// Если url не найден - отдаем соответствующий статус
	if errors.Is(err, shortenerErrors.ErrURLNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if errors.Is(err, shortenerErrors.ErrURLDeleted) {
		w.WriteHeader(http.StatusGone)
		w.Write([]byte("gone"))
		return
	}

	w.Header().Add("Location", url.URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte(url.URL))
}

// APICreateShortURLHTTPHandler Api для создания короткого урла
func APICreateShortURLHTTPHandler(w http.ResponseWriter, r *http.Request) {
	u := url{}

	// Обрабатываем входящий json
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uuid := middlewares.UserSignedCookie.UUID

	url, err := APICreateShortURLHandler(u.URL, uuid)

	// Если такой url уже есть - отдаем соответствующий статус
	if errors.Is(err, shortenerErrors.ErrURLConflict) {
		resp, _ := json.Marshal(response{URL: url.ShortURL})
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write(resp)
		return
	}

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

// APICreateShortURLBatchHTTPHandler Api для создания коротких урлов пачками
func APICreateShortURLBatchHTTPHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData []batchURL

	// Обрабатываем входящий json
	if err := json.NewDecoder(r.Body).Decode(&incomingData); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var urls []*types.URL
	var resp []*shortenBatchURL
	uuid := middlewares.UserSignedCookie.UUID

	for _, url := range incomingData {
		shortURL := fmt.Sprintf("%s/%s", app.Cfg.BaseURL, url.CorrelationID)

		urls = append(urls, &types.URL{
			UUID:     uuid,
			Hash:     url.CorrelationID,
			URL:      url.OriginalURL,
			ShortURL: shortURL,
		})
		resp = append(resp, &shortenBatchURL{
			CorrelationID: url.CorrelationID,
			ShortURL:      shortURL,
		})
	}

	_, err := APICreateShortURLBatchHandler(urls)
	if err != nil {
		fmt.Println(err)
	}

	response, _ := json.Marshal(resp)

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Accept", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

// APIDeleteShortURLBatchHTTPHandler удаляет урлы из базы по идентификаторам
func APIDeleteShortURLBatchHTTPHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData []string

	// Обрабатываем входящий json
	if err := json.NewDecoder(r.Body).Decode(&incomingData); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	APIDeleteShortURLBatchHandler(incomingData)

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("ok"))
}

// APIStatsHTTPHandler статистика по урлам
func APIStatsHTTPHandler(w http.ResponseWriter, r *http.Request) {
	resp := APIStatsHandler()

	respString, _ := json.Marshal(resp)

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Accept", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(respString)
}

// GetUserURLSHTTPHandler — возвращает все сокращенные урлы пользователя.
func GetUserURLSHTTPHandler(w http.ResponseWriter, r *http.Request) {
	uuid := middlewares.UserSignedCookie.UUID

	urls, err := GetUserURLSHandler(uuid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

// PingHTTPHandler проверяет соединение с базой
func PingHTTPHandler(w http.ResponseWriter, r *http.Request) {
	err := PingHandler()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("ok"))
}
