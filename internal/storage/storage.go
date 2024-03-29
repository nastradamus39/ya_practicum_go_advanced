package storage

import (
	"errors"
	"log"
	"os"

	shortenerErrors "github.com/nastradamus39/ya_practicum_go_advanced/internal/errors"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/types"
)

// Storage Хранилище ссылок
var Storage store

type repository interface {
	// Save сохраняет объект ссылки в хранилище
	Save(url *types.URL) error
	// FindByHash ищет урл в хранилище по хешу
	FindByHash(hash string) (exist bool, url *types.URL, err error)
	// FindByUUID ищет все ссылки пользователя с uuid
	FindByUUID(uuid string) (exist bool, urls map[string]*types.URL, err error)
	// DeleteByHash удаляет урлы
	DeleteByHash([]string) (err error)
}

type store interface {
	// Save сохраняет объект ссылки в хранилище
	Save(url *types.URL) error
	// SaveBatch сохраняет массив объектов ссылок в хранилище
	SaveBatch(urls []*types.URL) (err error)
	// FindByHash ищет урл в хранилище по хешу
	FindByHash(hash string) (exist bool, url *types.URL, err error)
	// FindByUUID ищет все ссылки пользователя с uuid
	FindByUUID(uuid string) (urls map[string]*types.URL, err error)
	// DeleteByHash удаляет урлы
	DeleteByHash([]string) (err error)
	// Drop чистит memory хранилище, удаляет файл
	Drop()
	// Ping Проверяет подключение к базе
	Ping() (err error)
	// Statistic Статистика
	Statistic() types.Statistic
}

type repositories struct {
	memory *MemoryRepository
	file   *FileRepository
	db     *DBRepository
}

type storage struct {
	cfg          *types.Config
	repositories repositories
}

func New(cfg *types.Config) (err error) {
	st := &storage{
		cfg: cfg,
	}

	mr := NewMemoryRepository()
	dbr := NewDBRepository(cfg)
	fr, err := NewFileRepository(cfg.DBPath)
	if err != nil {
		return err
	}

	// Инициируем репозитории
	st.repositories = repositories{
		memory: mr,
		file:   fr,
		db:     dbr,
	}

	Storage = st

	return nil
}

func (s *storage) Save(url *types.URL) (err error) {
	// Сохраняем в память
	err = s.repositories.memory.Save(url)
	// если не получилось записать в память - все плохо. выходим
	if err != nil {
		log.Println(err)
		return
	}

	// Сохраняем в файл
	if exist, _, _ := s.repositories.file.FindByHash(url.Hash); !exist {
		err = s.repositories.file.Save(url)
		// не получилось записать в файл - идем дальше
		if err != nil {
			log.Println(err)
		}
	}

	// Сохраняем в базу
	err = s.repositories.db.Save(url)
	// база опциональна
	if errors.Is(err, shortenerErrors.ErrNoDBConnection) {
		return nil
	}
	if err != nil {
		log.Println(err)
	}

	return
}

func (s *storage) SaveBatch(urls []*types.URL) (err error) {
	err = s.repositories.db.SaveBatch(urls)

	return
}

func (s *storage) DeleteByHash(urls []string) (err error) {
	err = s.repositories.db.DeleteByHash(urls)

	return
}

func (s *storage) FindByHash(hash string) (exist bool, url *types.URL, err error) {
	// Сначала в бд
	exist, url, err = s.repositories.db.FindByHash(hash)
	if exist {
		return
	}

	// ищем в файле
	exist, url, err = s.repositories.file.FindByHash(hash)
	if exist {
		return
	}

	// Ищем в памяти
	exist, url, err = s.repositories.memory.FindByHash(hash)
	if exist {
		return
	}

	return
}

func (s *storage) FindByUUID(uuid string) (urls map[string]*types.URL, err error) {
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

func (s *storage) Statistic() types.Statistic {
	stat := new(types.Statistic)

	stat.Urls = s.repositories.db.UrlsCount()
	stat.Users = s.repositories.db.UsersCount()

	return *stat
}

func (s *storage) Drop() {
	s.repositories.memory.items = map[string]*types.URL{}
	os.Remove(s.cfg.DBPath)
}

func (s *storage) Ping() (err error) {
	return s.repositories.db.Ping()
}
