package types

// Config конфиг приложения
type Config struct {
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	ServerPort    string `env:"SERVER_PORT" envDefault:"8080"`
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	DBPath        string `env:"FILE_STORAGE_PATH" envDefault:"./db"`
	DatabaseDsn   string `env:"DATABASE_DSN" envDefault:"postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable"`
}

// URL - структура для url
type URL struct {
	UUID     string
	Hash     string
	URL      string
	ShortURL string
}
