package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/rainset/shortener/internal/helper"
	"github.com/rainset/shortener/internal/storage"
	pb "github.com/rainset/shortener/proto"
)

// ShortenerServer поддерживает все необходимые методы сервера.
type ShortenerServer struct {
	config Config
	store  storage.InterfaceStorage
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedShortenerServer
}

// AddURL реализует интерфейс добавления ссылки.
func (s *ShortenerServer) AddURL(ctx context.Context, in *pb.AddURLRequest) (*pb.AddURLResponse, error) {
	var response pb.AddURLResponse
	hash := helper.GenerateToken(8)
	var isDBExist bool
	err := s.store.AddURL(hash, in.Url)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				isDBExist = true
				hash, err = s.store.GetByOriginalURL(in.Url)
				if err != nil {
					return &response, nil
				}
			}
		}
	}
	if err != nil && !isDBExist {
		response.Error = fmt.Sprintf("Ошибка при добавлении: %s", err)
		return &response, nil
	}
	response.Result = helper.GenerateShortenURL(s.config.ServerBaseURL, hash)

	return &response, nil
}

// GetURL реализует интерфейс получения ссылки по хешу.
func (s *ShortenerServer) GetURL(ctx context.Context, in *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	var response pb.GetURLResponse

	resultURL, err := s.store.GetURL(in.Hash)

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
func (s *ShortenerServer) Stats(ctx context.Context, in *pb.StatsRequest) (*pb.StatsResponse, error) {
	var response pb.StatsResponse

	stats, err := s.store.GetStats()
	if err != nil {
		response.Error = fmt.Sprintf("Ошибка получения данных: %s", err)
		return &response, nil
	}
	response.Urls = int64(stats.Urls)
	response.Users = int64(stats.Users)
	return &response, nil
}

// AddBatchURL реализует интерфейс массового добавления ссылок
func (s *ShortenerServer) AddBatchURL(ctx context.Context, in *pb.AddBatchURLRequest) (*pb.AddBatchURLResponse, error) {
	var response pb.AddBatchURLResponse

	batchURLs := make([]storage.BatchUrls, 0)
	for _, v := range in.Urls {
		batchURLs = append(batchURLs, storage.BatchUrls{CorrelationID: v.CorrelationId, OriginalURL: v.OriginalUrl})
	}
	result, err := s.store.AddBatchURL(batchURLs)
	if err != nil {
		response.Error = fmt.Sprintf("Ошибка массового добавления данных: %s", err)
		return &response, nil
	}

	urls := make([]*pb.BatchUrlResponse, 0)
	for _, v := range result {
		shortenURL := helper.GenerateShortenURL(s.config.ServerBaseURL, v.Hash)
		urls = append(urls, &pb.BatchUrlResponse{CorrelationId: v.CorrelationID, ShortUrl: shortenURL})
	}

	return &response, nil
}
