package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	handlers "github.com/nastradamus39/ya_practicum_go_advanced/internal/handlers"
	middlewares "github.com/nastradamus39/ya_practicum_go_advanced/internal/middleware"
	"log"
	"net/http"
)

type Config struct {
	BaseURL       string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
	ServerPort    string `env:"SERVER_PORT" envDefault:"8080"`
	ServerAddress string `env:"SERVER_HOST" envDefault:"localhost"`
	DbPath        string `env:"FILE_STORAGE_PATH" envDefault:"./db"`
}

var Cfg Config

func main() {
	r := router()

	err := env.Parse(&Cfg)
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&Cfg.ServerAddress, "server-host", Cfg.ServerAddress, "Базовый урл для коротких ссылок")
	flag.StringVar(&Cfg.ServerPort, "server-port", Cfg.ServerPort, "Порт сервера")
	flag.StringVar(&Cfg.BaseURL, "server-base", Cfg.BaseURL, "Адрес для запуска сервера")
	flag.StringVar(&Cfg.DbPath, "file-storage-path", Cfg.DbPath, "Путь к файлу с ссылками")
	flag.Parse()

	serverAddr := fmt.Sprintf("%s:%s", Cfg.ServerAddress, Cfg.ServerPort)
	bseUrl := fmt.Sprintf("http://%s:%s", Cfg.ServerAddress, Cfg.ServerPort)

	fmt.Println(fmt.Printf("Starting server on %s", serverAddr))
	fmt.Println(fmt.Printf("Base url %s", bseUrl))

	handlers.BaseUrl = bseUrl
	handlers.Storage, _ = handlers.NewFileStorage(Cfg.DbPath)

	http.ListenAndServe(serverAddr, r)
}

func router() (r *chi.Mux) {
	r = chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(middlewares.Decompress)

	r.Post("/", handlers.CreateShortURLHandler)
	r.Post("/api/shorten", handlers.ApiCreateShortURLHandler)
	r.Get("/{hash}", handlers.GetShortURLHandler)

	return r
}
