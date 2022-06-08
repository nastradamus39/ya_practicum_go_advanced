package app

import (
	"crypto/md5"
	"fmt"

	"github.com/nastradamus39/ya_practicum_go_advanced/internal/middlewares"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/storage"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/types"
)

// Cfg конфиг приложения
var Cfg types.Config

// Storage Хранилище ссылок в файле
var Storage *storage.FileStorage

// Urls ссылки в памяти
var Urls = map[string]map[string]*types.Url{}

// ShortURL сокращает переданный url, сохраняет, возвращает короткую ссылку
func ShortURL(originalUrl string) (url *types.Url) {
	h := md5.New()
	h.Write([]byte(originalUrl))

	hash := fmt.Sprintf("%x", h.Sum(nil))
	uuid := middlewares.UserSignedCookie.Uuid
	shortUrl := fmt.Sprintf("%s/%x", Cfg.BaseURL, h.Sum(nil))

	url = &types.Url{
		Uuid:     uuid,
		Hash:     hash,
		URL:      originalUrl,
		ShortUrl: shortUrl,
	}

	// Пишем в память
	if Urls[uuid] == nil {
		Urls[uuid] = map[string]*types.Url{}
	}
	Urls[uuid][hash] = url

	// Если есть в бд на диске - вернем. Если нет запишем
	if exist, _, _ := Storage.FindByHash(hash); !exist {
		Storage.Save(*url)
	}

	return url
}

// GetURLByHash возвращает полный url по хешу
func GetURLByHash(uuid string, hash string) (url *types.Url, err error) {
	// Ищем в памяти
	if url, exist := Urls[uuid][hash]; exist {
		return url, nil
	}

	// Если в памяти нет - ищем в файле
	if exist, u, _ := Storage.FindByHash(hash); exist {
		// Пишем в память
		if Urls[uuid] == nil {
			Urls[uuid] = map[string]*types.Url{}
		}
		Urls[uuid][hash] = &u
	}

	return Urls[uuid][hash], nil
}

// GetUrlsByUuid возвращает все ссылки по uuid
func GetUrlsByUuid(uuid string) (urls map[string]*types.Url) {
	// Ищем в памяти
	if urls, exist := Urls[uuid]; exist {
		return urls
	}

	// нет в памяти - ищем в файле
	if exist, urls, _ := Storage.FindByUuid(uuid); exist {
		// Пишем в память
		Urls[uuid] = urls
	}

	return Urls[uuid]
}
