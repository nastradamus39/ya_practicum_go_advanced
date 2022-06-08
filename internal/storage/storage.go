package storage

import (
	"encoding/json"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/types"
	"io"
	"os"
	"sync"
)

func newWriter(fileName string) (*writer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &writer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func newReader(fileName string) (*reader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &reader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func NewFileStorage(filename string) (fs *FileStorage, err error) {
	fs = &FileStorage{}
	fs.storageReader, err = newReader(filename)
	if err != nil {
		return nil, err
	}
	fs.storageWriter, err = newWriter(filename)
	if err != nil {
		return nil, err
	}
	return fs, nil
}

type writer struct {
	file    *os.File
	encoder *json.Encoder
}

func (p *writer) Write(url *types.Url) error {
	return p.encoder.Encode(&url)
}

func (p *writer) Close() error {
	return p.file.Close()
}

type reader struct {
	file    *os.File
	decoder *json.Decoder
}

func (c *reader) Read() (*types.Url, error) {
	event := &types.Url{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *reader) Close() error {
	return c.file.Close()
}

type FileStorage struct {
	mx            sync.Mutex
	storageReader *reader
	storageWriter *writer
}

// Save - сохраняет ID и ссылку в файле
func (f *FileStorage) Save(url types.Url) error {
	f.mx.Lock()
	defer f.mx.Unlock()

	err := f.storageWriter.Write(&url)
	if err != nil {
		return err
	}

	return nil
}

// FindByHash ищет в файле ссылку
func (f *FileStorage) FindByHash(hash string) (exist bool, url types.Url, err error) {
	f.mx.Lock()
	defer f.mx.Unlock()

	_, err = f.storageReader.file.Seek(0, io.SeekStart)
	if err != nil {
		return false, types.Url{}, err
	}

	for {
		url, err := f.storageReader.Read()

		if err != nil {
			return false, types.Url{}, err
		}

		if url.Hash == hash {
			return true, *url, nil
		}
	}
}

// FindByUuid ищет в файле ссылки по uuid
func (f *FileStorage) FindByUuid(uuid string) (exist bool, urls map[string]*types.Url, err error) {
	f.mx.Lock()
	defer f.mx.Unlock()

	_, err = f.storageReader.file.Seek(0, io.SeekStart)
	if err != nil {
		return false, map[string]*types.Url{}, err
	}

	urls = map[string]*types.Url{}

	for {
		url, err := f.storageReader.Read()

		if err != nil {
			break
		}

		if url.Uuid == uuid {
			urls[url.Hash] = url
		}
	}

	return true, urls, nil
}
