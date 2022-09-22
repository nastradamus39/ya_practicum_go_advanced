package types

import "database/sql"

// Config конфиг приложения
type Config struct {
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080" json:"base_url"`
	ServerPort    string `env:"SERVER_PORT" envDefault:"8080"`
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080" json:"server_address"`
	DBPath        string `env:"FILE_STORAGE_PATH" envDefault:"./db" json:"file_storage_path"`
	DatabaseDsn   string `env:"DATABASE_DSN" envDefault:"" json:"database_dsn"`
	EnableHttps   bool   `env:"ENABLE_HTTPS" envDefault:"true" json:"enable_https"`
}

// URL - структура для url
type URL struct {
	UUID      string         `db:"uuid"`
	Hash      string         `db:"hash"`
	URL       string         `db:"url"`
	ShortURL  string         `db:"short_url"`
	DeletedAt sql.NullString `db:"deleted_at"`
}

// Statistic - статистика
type Statistic struct {
	Urls  int `json:"urls"`
	Users int `json:"users"`
}
