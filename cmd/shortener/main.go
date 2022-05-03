package main

import (
	"crypto/md5"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io/ioutil"
	"net/http"
)

var urls = map[string]string{}

// PostUrl — создает короткий урл.
func PostUrl(w http.ResponseWriter, r *http.Request) {
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

// GetShortUrl — возвращает полный урл по короткому.
func GetShortUrl(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	url := urls[hash]

	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)

	w.Write([]byte(url))
}

func main() {
	r := router()

	fmt.Println("Starting server on :8080")
	http.ListenAndServe("127.0.0.1:8080", r)
}

func router() (r *chi.Mux) {
	r = chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", PostUrl)
	r.Get("/{hash}", GetShortUrl)

	return r
}
