package storage

import (
	"crypto/md5"
	"fmt"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/types"
)

type MemoryRepository struct {
	items map[string]*types.URL
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		items: map[string]*types.URL{},
	}
}

func (r *MemoryRepository) Save(url *types.URL) error {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s_%s", url.UUID, url.URL)))

	key := fmt.Sprintf("%x", h.Sum(nil))

	r.items[key] = url

	return nil
}

func (r *MemoryRepository) FindByHash(hash string) (exist bool, url *types.URL, err error) {
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

func (r *MemoryRepository) FindByUUID(uuid string) (urls map[string]*types.URL, err error) {
	urls = map[string]*types.URL{}
	err = nil

	for _, item := range r.items {
		if item.UUID == uuid {
			urls[item.Hash] = item
		}
	}

	return
}
