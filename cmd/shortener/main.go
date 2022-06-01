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
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	ServerPort    string `env:"SERVER_PORT" envDefault:"8080"`
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	DBPath        string `env:"FILE_STORAGE_PATH" envDefault:"./db"`
}

var Cfg Config

func main() {
	r := router()

	err := env.Parse(&Cfg)
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&Cfg.ServerAddress, "a", Cfg.ServerAddress, "Адрес для запуска сервера")
	flag.StringVar(&Cfg.ServerPort, "server-port", Cfg.ServerPort, "Порт сервера")
	flag.StringVar(&Cfg.BaseURL, "b", Cfg.BaseURL, "Базовый адрес результирующего сокращённого URL")
	flag.StringVar(&Cfg.DBPath, "f", Cfg.DBPath, "Путь к файлу с ссылками")
	flag.Parse()

	serverAddr := Cfg.ServerAddress
	baseURL := Cfg.BaseURL

	fmt.Println(fmt.Printf("Starting server on %s", serverAddr))
	fmt.Println(fmt.Printf("Base url %s", baseURL))

	handlers.BaseURL = baseURL
	handlers.Storage, _ = handlers.NewFileStorage(Cfg.DBPath)

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
	r.Use(middlewares.UserCookie)

	r.Post("/", handlers.CreateShortURLHandler)
	r.Post("/api/shorten", handlers.APICreateShortURLHandler)
	r.Get("/{hash}", handlers.GetShortURLHandler)
	r.Get("/api/user/urls", handlers.GetShortURLHandler)

	return r
}
