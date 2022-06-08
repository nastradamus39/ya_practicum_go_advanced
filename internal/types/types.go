package types

// Config конфиг приложения
type Config struct {
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	ServerPort    string `env:"SERVER_PORT" envDefault:"8080"`
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	DBPath        string `env:"FILE_STORAGE_PATH" envDefault:"./db"`
}

// Url - структура для url
type Url struct {
	Uuid     string
	Hash     string
	URL      string
	ShortUrl string
}
