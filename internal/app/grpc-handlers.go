package app

import (
	"context"
	"fmt"
	"github.com/rainset/shortener/internal/proto"
)

// ShortenerServer поддерживает все необходимые методы сервера.
type ShortenerGRPCServer struct {
	a *App
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	proto.UnimplementedShortenerServer
}

// AddURL реализует интерфейс добавления ссылки.
func (s *ShortenerGRPCServer) AddURL(ctx context.Context, in *proto.AddURLRequest) (*proto.AddURLResponse, error) {

	var response proto.AddURLResponse
	addURLResult, err := s.a.AddURL(in.Url)
	if err != nil {
		response.Error = fmt.Sprintf("Ошибка при добавлении: %s", err)
		return &response, nil
	}
	response.Result = addURLResult.ShortURL

	return &response, nil
}

// GetURL реализует интерфейс получения ссылки по хешу.
func (s *ShortenerGRPCServer) GetURL(ctx context.Context, in *proto.GetURLRequest) (*proto.GetURLResponse, error) {
	var response proto.GetURLResponse

	resultURL, err := s.a.GetURL(in.Hash)

	if err != nil {
		response.Error = fmt.Sprintf("Ошибка получения данных: %s", err)
		return &response, nil
	}

	if resultURL.Deleted == 1 {
		response.Error = "Ссылка удалена"
		return &response, nil
	}
	response.Result = resultURL.Original

	return &response, nil
}

// Stats реализует интерфейс получения статистики
func (s *ShortenerGRPCServer) Stats(ctx context.Context, in *proto.StatsRequest) (*proto.StatsResponse, error) {
	var response proto.StatsResponse

	stats, err := s.a.GetStats()
	if err != nil {
		response.Error = fmt.Sprintf("Ошибка получения данных: %s", err)
		return &response, nil
	}
	response.Urls = int64(stats.Urls)
	response.Users = int64(stats.Users)
	return &response, nil
}

// AddBatchURL реализует интерфейс массового добавления ссылок
func (s *ShortenerGRPCServer) AddBatchURL(ctx context.Context, in *proto.AddBatchURLRequest) (*proto.AddBatchURLResponse, error) {
	var response proto.AddBatchURLResponse

	batchURLs := make([]AddURLBatchRequest, 0)
	for _, v := range in.Urls {
		batchURLs = append(batchURLs, AddURLBatchRequest{CorrelationID: v.Correlation_ID, OriginalURL: v.OriginalUrl})
	}

	result, err := s.a.AddBatchURL(batchURLs)
	if err != nil {
		response.Error = fmt.Sprintf("Ошибка массового добавления данных: %s", err)
		return &response, nil
	}

	urls := make([]*proto.BatchUrlResponse, 0)
	for _, v := range result {
		urls = append(urls, &proto.BatchUrlResponse{Correlation_ID: v.CorrelationID, ShortUrl: v.ShortURL})
	}
	response.Urls = urls

	return &response, nil
}
