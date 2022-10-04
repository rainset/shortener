package memory

import (
	"fmt"
	"github.com/rainset/shortener/internal/storage"
	"sync"
)

type Memory struct {
	mutex           sync.RWMutex
	urls            map[string]ResultURL
	userHistoryURLs []UserHistoryURL
}

type ResultURL struct {
	Hash     string
	Original string
}

type UserHistoryURL struct {
	CookieID string
	Hash     string
}

func New() *Memory {
	urls := make(map[string]ResultURL, 0)
	userHistoryURLs := make([]UserHistoryURL, 0)
	return &Memory{
		urls:            urls,
		userHistoryURLs: userHistoryURLs,
	}
}

func (m *Memory) AddURL(hash, original string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.urls[hash] = ResultURL{Hash: hash, Original: original}
	fmt.Println("urls", m.urls)
	return nil
}

func (m *Memory) GetURL(hash string) (original string, err error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	item := m.urls[hash]
	return item.Original, nil
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

func (m *Memory) DeleteBatchURL(_ string, _ []string) (err error) {
	return err
}

func (m *Memory) AddUserHistoryURL(cookieID, hash string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.userHistoryURLs = append(m.userHistoryURLs, UserHistoryURL{Hash: hash, CookieID: cookieID})

	fmt.Println("userHistoryURLs", m.userHistoryURLs)
	return nil
}

func (m *Memory) GetListUserHistoryURL(cookieID string) (result []storage.ResultHistoryURL, err error) {
	for _, item := range m.userHistoryURLs {
		if cookieID == item.CookieID {
			original := m.urls[item.Hash].Original
			result = append(result, storage.ResultHistoryURL{CookieID: cookieID, Hash: item.Hash, Original: original})
		}
	}
	fmt.Println("GetListUserHistoryURL", result)
	return result, err
}

func (m *Memory) Ping() (err error) {
	return err
}
