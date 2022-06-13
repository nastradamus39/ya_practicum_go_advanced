package utils

import (
	"crypto/md5"
	"fmt"

	"github.com/nastradamus39/ya_practicum_go_advanced/internal/app"
)

// GetShortUrl создает короткий урл из полного и возвращает хеш
func GetShortUrl(value string) (hash string, shortUrl string) {
	h := md5.New()
	h.Write([]byte(value))

	hash = fmt.Sprintf("%x", h.Sum(nil))
	shortUrl = fmt.Sprintf("%s/%x", app.Cfg.BaseURL, h.Sum(nil))

	return
}
