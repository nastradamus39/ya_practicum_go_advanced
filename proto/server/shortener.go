package server

import (
	"context"
	"fmt"

	"github.com/nastradamus39/ya_practicum_go_advanced/internal/handlers"
	proto "github.com/nastradamus39/ya_practicum_go_advanced/proto"
)

// ShortenerServer поддерживает все необходимые методы сервера.
type ShortenerServer struct {
	// Нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	proto.UnimplementedUrlsServer
}

// CreateShortURLHandler реализует интерфейс добавления пользователя.
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

// GetShortURLHandler реализует интерфейс добавления пользователя.
func (s *ShortenerServer) GetShortURLHandler(ctx context.Context, in *proto.GetUrlRequest) (*proto.GetUrlResponse, error) {
	fmt.Println("GetShortURLHandler")
	return nil, nil

	var response proto.GetUrlResponse

	url, err := handlers.GetShortURLHandler(in.Hash)

	if err == nil {
		response.Url = url.ShortURL
		return &response, nil
	}

	return &response, nil
}
