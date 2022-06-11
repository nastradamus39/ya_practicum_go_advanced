package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/app"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/handlers"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/middlewares"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/storage"
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := Router()

	// Переменные окружения в конфиг
	err := env.Parse(&app.Cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Параметры командной строки в конфиг
	flag.StringVar(&app.Cfg.ServerAddress, "a", app.Cfg.ServerAddress, "Адрес для запуска сервера")
	flag.StringVar(&app.Cfg.ServerPort, "server-port", app.Cfg.ServerPort, "Порт сервера")
	flag.StringVar(&app.Cfg.BaseURL, "b", app.Cfg.BaseURL, "Базовый адрес результирующего сокращённого URL")
	flag.StringVar(&app.Cfg.DBPath, "f", app.Cfg.DBPath, "Путь к файлу с ссылками")
	flag.StringVar(&app.Cfg.DatabaseDsn, "d", app.Cfg.DatabaseDsn, "Строка с адресом подключения к БД")
	flag.Parse()

	fmt.Println(fmt.Printf("Starting server on %s", app.Cfg.ServerAddress))
	fmt.Println(fmt.Printf("Base url %s", app.Cfg.BaseURL))

	// инициируем хранилище
	s := storage.Storage{}
	app.Storage, _ = s.New(&app.Cfg)

	// подключение к бд
	app.DB, err = sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		fmt.Println(err)
	}
	defer app.DB.Close()

	// запускаем сервер
	http.ListenAndServe(app.Cfg.ServerAddress, r)
}

func Router() (r *chi.Mux) {
	r = chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(middlewares.Decompress)
	r.Use(middlewares.UserCookie)

	r.Post("/", handlers.CreateShortURLHandler)
	r.Get("/ping", handlers.PingHandler)
	r.Get("/api/user/urls", handlers.GetUserURLSHandler)
	r.Post("/api/shorten", handlers.APICreateShortURLHandler)
	r.Get("/{hash}", handlers.GetShortURLHandler)

	return r
}
