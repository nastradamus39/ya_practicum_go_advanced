package types

// Config конфиг приложения
type Config struct {
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	ServerPort    string `env:"SERVER_PORT" envDefault:"8080"`
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	DBPath        string `env:"FILE_STORAGE_PATH" envDefault:"./db"`
	DatabaseDsn   string `env:"DATABASE_DSN" envDefault:""`
}

// URL - структура для url
type URL struct {
	UUID     string `db:"uuid"`
	Hash     string `db:"hash"`
	URL      string `db:"url"`
	ShortURL string `db:"short_url"`
}
