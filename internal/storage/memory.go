package storage

import (
	"crypto/md5"
	"fmt"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/types"
)

type MemoryRepository struct {
	items map[string]*types.Url
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		items: map[string]*types.Url{},
	}
}

func (r *MemoryRepository) Save(url *types.Url) error {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s_%s", url.Uuid, url.URL)))

	key := fmt.Sprintf("%x", h.Sum(nil))

	r.items[key] = url

	return nil
}

func (r *MemoryRepository) FindByHash(hash string) (exist bool, url *types.Url, err error) {
	exist = false
	url = nil
	err = nil

	for _, item := range r.items {
		if item.Hash == hash {
			url = item
			exist = true
		}
	}

	return
}

func (r *MemoryRepository) FindByUuid(uuid string) (urls map[string]*types.Url, err error) {
	urls = map[string]*types.Url{}
	err = nil

	for _, item := range r.items {
		if item.Uuid == uuid {
			urls[item.Hash] = item
		}
	}

	return
}
