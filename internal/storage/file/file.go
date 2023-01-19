// Пакет для работы с базой данных размещенной в файловой системе
package file

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"

	"github.com/rainset/shortener/internal/storage"
)

// File -
type File struct {
	mutex           sync.Mutex
	fileStoragePath string
	userHistoryURLs []UserHistoryURL
}

// ResultURL -
type ResultURL struct {
	Hash     string
	Original string
}

// UserHistoryURL -
type UserHistoryURL struct {
	CookieID string
	Hash     string
}

// DataURL -
type DataURL struct {
	Hash    string `json:"hash"`
	LongURL string `json:"long_url"`
}

// producer -
type producer struct {
	file    *os.File
	encoder *json.Encoder
}

// consumer -
type consumer struct {
	file    *os.File
	decoder *json.Decoder
}

// New -
func New(fileStoragePath string) *File {

	var userHistoryURLs []UserHistoryURL
	return &File{
		fileStoragePath: fileStoragePath,
		userHistoryURLs: userHistoryURLs,
	}
}

// NewProducer -
func NewProducer(fileName string) (*producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

// WriteURL -
func (p *producer) WriteURL(url *DataURL) error {
	return p.encoder.Encode(&url)
}

// Close -
func (p *producer) Close() error {
	return p.file.Close()
}

// NewConsumer -
func NewConsumer(fileName string) (*consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

// ReadURL -
func (c *consumer) ReadURL() (*DataURL, error) {
	url := &DataURL{}
	if err := c.decoder.Decode(&url); err != nil {
		return nil, err
	}
	return url, nil
}

// Close -
func (c *consumer) Close() error {
	return c.file.Close()
}

// RestoreStorage -
func (c *consumer) RestoreStorage() (result []ResultURL, err error) {
	for {
		dataURL := &DataURL{}
		if err := c.decoder.Decode(&dataURL); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		result = append(result, ResultURL{Hash: dataURL.Hash, Original: dataURL.LongURL})
	}
	return result, nil
}

// AddURL -
func (f *File) AddURL(hash, original string) (err error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	producer, err := NewProducer(f.fileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	requestData := &DataURL{Hash: hash, LongURL: original}
	if err = producer.WriteURL(requestData); err != nil {
		return err
	}
	defer producer.Close()
	return nil
}

// GetURL -
func (f *File) GetURL(hash string) (resultURL storage.ResultURL, err error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	consumer, err := NewConsumer(f.fileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	urls, err := consumer.RestoreStorage()
	for _, item := range urls {
		if item.Hash == hash {
			resultURL.Original = item.Original
			return resultURL, nil
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	return resultURL, err
}

// GetByOriginalURL -
func (f *File) GetByOriginalURL(original string) (hash string, err error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	consumer, err := NewConsumer(f.fileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	urls, err := consumer.RestoreStorage()
	for _, item := range urls {
		if item.Original == original {
			return item.Hash, nil
		}
	}
	return "", err
}

// AddBatchURL заглушка реализовано только для postgres
func (f *File) AddBatchURL(_ []storage.BatchUrls) (result []storage.ResultBatchUrls, err error) {
	return result, err
}

// DeleteUserBatchURL заглушка реализовано только для postgres
func (f *File) DeleteUserBatchURL(_ string, _ []string) (err error) {
	return err
}

// DeleteBatchURL заглушка реализовано только для postgres
func (f *File) DeleteBatchURL(_ []string) (err error) {
	return err
}

// AddUserHistoryURL -
func (f *File) AddUserHistoryURL(cookieID, hash string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.userHistoryURLs = append(f.userHistoryURLs, UserHistoryURL{Hash: hash, CookieID: cookieID})
	return nil
}

// GetListUserHistoryURL -
func (f *File) GetListUserHistoryURL(cookieID string) (result []storage.ResultHistoryURL, err error) {

	urlList := make(map[string]string)
	consumer, err := NewConsumer(f.fileStoragePath)
	if err != nil {
		log.Fatal(err)
	}

	urls, err := consumer.RestoreStorage()
	for _, item := range urls {
		urlList[item.Hash] = item.Original
	}

	for _, row := range f.userHistoryURLs {
		if cookieID == row.CookieID {
			original := urlList[row.Hash]
			result = append(result, storage.ResultHistoryURL{CookieID: cookieID, Hash: row.Hash, Original: original})
		}
	}
	return result, err
}

// Ping -
func (f *File) Ping() (err error) {
	return err
}

// GetStats -
func (f *File) GetStats() (stats storage.Stats, err error) {
	return stats, err
}
