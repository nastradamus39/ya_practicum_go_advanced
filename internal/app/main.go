package app

import (
	"database/sql"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/storage"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/types"
)

// Cfg конфиг приложения
var Cfg types.Config

// Storage Хранилище ссылок
var Storage *storage.Storage

// DB база
var DB *sql.DB
