package storage

import (
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/types"
	"os"
)

type repository interface {
	// Save сохраняет объект ссылки в хранилище
	Save(url *types.URL) error
	// FindByHash ищет урл в хранилище по хешу
	FindByHash(hash string) (exist bool, url *types.URL, err error)
	// FindByUuid ищет все ссылки пользователя с uuid
	FindByUuid(uuid string) (exist bool, urls map[string]*types.URL, err error)
}

type storage interface {
	// New инициирует хранилище
	New(cfg *types.Config) (*Storage, error)
	// Save сохраняет объект ссылки в хранилище
	Save(url *types.URL) error
	// FindByHash ищет урл в хранилище по хешу
	FindByHash(hash string) (exist bool, url *types.URL, err error)
	// FindByUuid ищет все ссылки пользователя с uuid
	FindByUuid(uuid string) (urls map[string]*types.URL, err error)
	// Drop чистит memory хранилище, удаляет файл
	Drop()
}

type repositories struct {
	memory *MemoryRepository
	file   *FileRepository
}

type Storage struct {
	cfg          *types.Config
	repositories repositories
}

func (s *Storage) New(cfg *types.Config) (*Storage, error) {
	s = &Storage{
		cfg: cfg,
	}

	mr := NewMemoryRepository()
	fr, err := NewFileRepository(cfg.DBPath)
	if err != nil {
		return nil, err
	}

	// Инициируем репозитории
	s.repositories = repositories{
		memory: mr,
		file:   fr,
	}

	return s, nil
}

func (s *Storage) Save(url *types.URL) error {
	// Сохраняем в память
	err := s.repositories.memory.Save(url)
	if err != nil {
		return err
	}

	// Сохраняем в файл
	if exist, _, _ := s.repositories.file.FindByHash(url.Hash); !exist {
		return s.repositories.file.Save(url)
	}

	return nil
}

func (s *Storage) FindByHash(hash string) (exist bool, url *types.URL, err error) {
	// Ищем в памяти
	exist, url, err = s.repositories.memory.FindByHash(hash)

	// Если есть в памяти - дальше не ищем
	if exist {
		return
	}

	// ищем в файле
	exist, url, err = s.repositories.file.FindByHash(hash)

	return
}

func (s *Storage) FindByUUID(uuid string) (urls map[string]*types.URL, err error) {
	// Ищем в памяти
	um, e := s.repositories.memory.FindByUUID(uuid)
	if e != nil {
		return nil, e
	}

	// Ищем в файле
	uf, e := s.repositories.file.FindByUUID(uuid)
	if e != nil {
		return nil, e
	}

	urls = map[string]*types.URL{}
	for _, item := range um {
		urls[item.Hash] = item
	}
	for _, item := range uf {
		urls[item.Hash] = item
	}

	return urls, nil
}

func (s *Storage) Drop() {
	s.repositories.memory.items = map[string]*types.URL{}
	os.Remove(s.cfg.DBPath)
}
