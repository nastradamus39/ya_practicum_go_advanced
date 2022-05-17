package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	handlers "github.com/nastradamus39/ya_practicum_go_advanced/internal/handlers"
)

func main() {
	r := router()

	// export SERVER_ADDRESS='127.0.0.1'
	// export BASE_URL='127.0.0.1'
	//os.Setenv("SERVER_ADDRESS", "127.0.0.1:8080")
	//os.Setenv("BASE_URL", "http://127.0.0.1:8080")

	//os.Setenv("SERVER_HOST", "127.0.0.1")
	//os.Setenv("SERVER_PORT", "8080")

	serverAddr := fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))

	os.Setenv("BASE_URL", fmt.Sprintf("http://%s", serverAddr))

	fmt.Println(fmt.Printf("Starting server on %s", serverAddr))

	http.ListenAndServe(serverAddr, r)
}

func router() (r *chi.Mux) {
	r = chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", handlers.CreateShortURLHandler)
	r.Post("/api/shorten", handlers.ApiCreateShortURLHandler)
	r.Get("/{hash}", handlers.GetShortURLHandler)

	return r
}
