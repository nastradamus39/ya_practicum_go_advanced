package handlers

import (
	"errors"
	"fmt"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/app"
	shortenerErrors "github.com/nastradamus39/ya_practicum_go_advanced/internal/errors"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/storage"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/types"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/utils"
	"log"
)

// CreateShortURLHandler — создает короткий урл.
func CreateShortURLHandler(originalURL string, uuid string) (url *types.URL, err error) {
	hash, shortURL := utils.GetShortURL(originalURL)

	url = &types.URL{
		UUID:     uuid,
		Hash:     hash,
		URL:      originalURL,
		ShortURL: shortURL,
	}

	err = storage.Storage.Save(url)

	// Если такой url уже есть - отдаем соответствующий статус
	if errors.Is(err, shortenerErrors.ErrURLConflict) {
		return url, err
	}

	// Другие ошибки при сохранении в хранилище
	if err != nil {
		log.Printf("CreateShortURLHandler. Не удалось сохранить урл в хранилище. %s", err)
		return nil, err
	}

	return url, nil
}

// GetShortURLHandler — возвращает полный урл по короткому.
func GetShortURLHandler(hash string) (url *types.URL, err error) {
	var exist bool

	exist, url, err = storage.Storage.FindByHash(hash)

	if !exist {
		return nil, shortenerErrors.ErrURLNotFound
	}

	if err != nil {
		return nil, err
	}

	deletedAt, _ := url.DeletedAt.Value()

	if deletedAt != nil {
		return nil, shortenerErrors.ErrURLDeleted
	}

	return url, nil
}

// APICreateShortURLHandler Api для создания короткого урла
func APICreateShortURLHandler(originalURL string, uuid string) (url *types.URL, err error) {
	hash, shortURL := utils.GetShortURL(originalURL)

	url = &types.URL{
		UUID:     uuid,
		Hash:     hash,
		URL:      originalURL,
		ShortURL: shortURL,
	}

	err = storage.Storage.Save(url)

	if err != nil {
		return url, err
	}

	return url, nil
}

// APICreateShortURLBatchHandler Api для создания коротких урлов пачками
func APICreateShortURLBatchHandler(urls []*types.URL) ([]*types.URL, error) {
	// Вычисляем короткий url для каждой ссылки
	for _, url := range urls {
		url.ShortURL = fmt.Sprintf("%s/%s", app.Cfg.BaseURL, url.Hash)
	}

	err := storage.Storage.SaveBatch(urls)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return urls, nil
}

// APIDeleteShortURLBatchHandler удаляет урлы из базы по идентификаторам
func APIDeleteShortURLBatchHandler(hashes []string) {
	if len(hashes) > 0 {
		go storage.Storage.DeleteByHash(hashes)
	}
}

// APIStatsHandler статистика по урлам
func APIStatsHandler() (statistic types.Statistic) {
	return storage.Storage.Statistic()
}

// GetUserURLSHandler — возвращает все сокращенные урлы пользователя.
func GetUserURLSHandler(uuid string) (urls map[string]*types.URL, err error) {
	urls, err = storage.Storage.FindByUUID(uuid)

	return urls, err
}

// PingHandler проверяет соединение с базой
func PingHandler() (err error) {
	return storage.Storage.Ping()
}
