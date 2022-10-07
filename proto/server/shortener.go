package server

import (
	"context"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/types"
	"log"

	"github.com/nastradamus39/ya_practicum_go_advanced/internal/handlers"
	proto "github.com/nastradamus39/ya_practicum_go_advanced/proto"
)

// ShortenerServer поддерживает все необходимые методы сервера.
type ShortenerServer struct {
	// Нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	proto.UnimplementedUrlsServer
}

// CreateShortURLHandler создает короткий url
func (s *ShortenerServer) CreateShortURLHandler(ctx context.Context, in *proto.AddUrlRequest) (*proto.AddUrlResponse, error) {
	var response proto.AddUrlResponse

	url, err := handlers.CreateShortURLHandler(in.Url, in.Uuid)

	if err == nil {
		response.Url = url.ShortURL
		return &response, nil
	}

	response.Error = err.Error()

	return &response, nil
}

// GetShortURLHandler возвращает короткий url
func (s *ShortenerServer) GetShortURLHandler(ctx context.Context, in *proto.GetUrlRequest) (*proto.GetUrlResponse, error) {
	var response proto.GetUrlResponse

	url, err := handlers.GetShortURLHandler(in.Hash)

	if err == nil {
		response.Url = url.ShortURL
		return &response, nil
	}

	return &response, nil
}

// APICreateShortURLHandler возвращает короткий url
func (s *ShortenerServer) APICreateShortURLHandler(ctx context.Context, in *proto.APICreateShortURLRequest) (*proto.APICreateShortURLResponse, error) {
	var response proto.APICreateShortURLResponse

	url, err := handlers.GetShortURLHandler(in.OriginalURL)

	if err == nil {
		response.ShortURL = url.ShortURL
		response.URL = url.URL
		response.Hash = url.Hash
		return &response, nil
	}

	return &response, nil
}

// APICreateShortURLBatchHandler апи пакетного создания url
func (s *ShortenerServer) APICreateShortURLBatchHandler(ctx context.Context, in *proto.APICreateShortURLBatchRequest) (*proto.APICreateShortURLBatchResponse, error) {
	var response proto.APICreateShortURLBatchResponse
	var urls []*types.URL

	for _, url := range in.Urls {
		u, e := handlers.GetShortURLHandler(url)

		if e == nil {
			urls = append(urls, u)
		} else {
			log.Println(e)
		}
	}

	//r, e := handlers.APICreateShortURLBatchHandler(urls)

	return nil, nil
}
