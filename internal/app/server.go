package app

import (
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/rainset/shortener/internal/helper"
	"github.com/rainset/shortener/internal/storage"
)

type ListURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type AddURLResult struct {
	ShortURL string
	Hash     string
	Exists   bool
}

type AddURLBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type AddURLBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (a *App) GetStats() (storage.Stats, error) {

	stats, err := a.s.GetStats()
	return stats, err
}

func (a *App) Ping() error {

	err := a.s.Ping()
	return err
}

func (a *App) GetURL(hash string) (storage.ResultURL, error) {

	result, err := a.s.GetURL(hash)
	return result, err
}

func (a *App) GetListUserHistoryURL(userID string) ([]ListURL, error) {

	var userHistoryURLs []storage.ResultHistoryURL
	userHistoryURLs, err := a.s.GetListUserHistoryURL(userID)

	list := make([]ListURL, 0)

	for _, item := range userHistoryURLs {
		if len(item.Hash) == 0 || len(item.CookieID) == 0 {
			continue
		}
		ShortURL := helper.GenerateShortenURL(a.Config.ServerBaseURL, item.Hash)
		list = append(list, ListURL{ShortURL: ShortURL, OriginalURL: item.Original})
	}

	return list, err
}

func (a *App) AddURL(originalURL string) (result AddURLResult, err error) {

	result.Hash = helper.GenerateToken(8)
	result.Exists = false

	err = a.s.AddURL(result.Hash, originalURL)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				result.Hash, err = a.s.GetByOriginalURL(originalURL)
				if err != nil {
					return result, err
				}
				result.Exists = true
			}
		}
	}
	result.ShortURL = helper.GenerateShortenURL(a.Config.ServerBaseURL, result.Hash)
	return result, err
}

func (a *App) AddBatchURL(list []AddURLBatchRequest) (result []AddURLBatchResponse, err error) {

	batchURLs := make([]storage.BatchUrls, 0)
	for _, v := range list {
		batchURLs = append(batchURLs, storage.BatchUrls{CorrelationID: v.CorrelationID, OriginalURL: v.OriginalURL})
	}

	addResult, err := a.s.AddBatchURL(batchURLs)

	if err != nil {
		return result, err
	}

	for _, v := range addResult {
		shortenURL := helper.GenerateShortenURL(a.Config.ServerBaseURL, v.Hash)
		result = append(result, AddURLBatchResponse{ShortURL: shortenURL, CorrelationID: v.CorrelationID})
	}

	return result, err
}

func (a *App) AddUserHistoryURL(userId, hash string) (err error) {
	err = a.s.AddUserHistoryURL(userId, hash)
	return err
}