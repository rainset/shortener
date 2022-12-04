// Пакет для работы с базой данных размещенной в памяти приложения
package memory

import (
	"sync"

	"github.com/rainset/shortener/internal/storage"
)

type Memory struct {
	mutex           sync.RWMutex
	urls            map[string]storage.ResultURL
	userHistoryURLs []UserHistoryURL
}

type UserHistoryURL struct {
	CookieID string
	Hash     string
}

func New() *Memory {
	urls := make(map[string]storage.ResultURL, 0)
	userHistoryURLs := make([]UserHistoryURL, 0)
	return &Memory{
		urls:            urls,
		userHistoryURLs: userHistoryURLs,
	}
}

func (m *Memory) AddURL(hash, original string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.urls[hash] = storage.ResultURL{ID: 0, Hash: hash, Original: original, Deleted: 0}
	return nil
}

func (m *Memory) GetURL(hash string) (resultURL storage.ResultURL, err error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	resultURL = m.urls[hash]
	return resultURL, nil
}

func (m *Memory) GetByOriginalURL(original string) (hash string, err error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for _, item := range m.urls {
		if item.Original == original {
			return item.Hash, nil
		}
	}
	return "", err
}

//заглушка реализовано только для postgres

func (m *Memory) AddBatchURL(_ []storage.BatchUrls) (result []storage.ResultBatchUrls, err error) {
	return result, err
}

func (m *Memory) DeleteUserBatchURL(cookieID string, hashes []string) (err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, hash := range hashes {
		for _, uh := range m.userHistoryURLs {
			if cookieID == uh.CookieID && hash == uh.Hash {
				m.urls[hash] = storage.ResultURL{ID: 0, Hash: hash, Original: m.urls[hash].Original, Deleted: 1}
			}
		}
	}
	return err
}

func (m *Memory) DeleteBatchURL(hashes []string) (err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, hash := range hashes {
		m.urls[hash] = storage.ResultURL{ID: 0, Hash: hash, Original: m.urls[hash].Original, Deleted: 1}
	}
	return err
}

func (m *Memory) AddUserHistoryURL(cookieID, hash string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.userHistoryURLs = append(m.userHistoryURLs, UserHistoryURL{Hash: hash, CookieID: cookieID})
	return nil
}

func (m *Memory) GetListUserHistoryURL(cookieID string) (result []storage.ResultHistoryURL, err error) {
	for _, item := range m.userHistoryURLs {
		if cookieID == item.CookieID {
			original := m.urls[item.Hash].Original
			result = append(result, storage.ResultHistoryURL{CookieID: cookieID, Hash: item.Hash, Original: original})
		}
	}
	return result, err
}

func (m *Memory) Ping() (err error) {
	return err
}
