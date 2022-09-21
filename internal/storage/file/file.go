package file

import (
	"encoding/json"
	"github.com/rainset/shortener/internal/storage"
	"io"
	"log"
	"os"
	"sync"
)

type File struct {
	mutex           sync.Mutex
	fileStoragePath string
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

type DataURL struct {
	Hash    string `json:"hash"`
	LongURL string `json:"long_url"`
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

type consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func New(fileStoragePath string) *File {

	var userHistoryURLs []UserHistoryURL
	return &File{
		fileStoragePath: fileStoragePath,
		userHistoryURLs: userHistoryURLs,
	}
}

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

func (p *producer) WriteURL(url *DataURL) error {
	return p.encoder.Encode(&url)
}

func (p *producer) Close() error {
	return p.file.Close()
}

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

func (c *consumer) ReadURL() (*DataURL, error) {
	url := &DataURL{}
	if err := c.decoder.Decode(&url); err != nil {
		return nil, err
	}
	return url, nil
}

func (c *consumer) Close() error {
	return c.file.Close()
}

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

func (f *File) GetURL(hash string) (original string, err error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	consumer, err := NewConsumer(f.fileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	urls, err := consumer.RestoreStorage()
	for _, item := range urls {
		if item.Hash == hash {
			return item.Original, nil
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	return "", err
}

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

// заглушка реализовано только для postgres

func (f *File) AddBatchURL(_ []storage.BatchUrls) (result []storage.ResultBatchUrls, err error) {
	return result, err
}

func (f *File) AddUserHistoryURL(cookieID, hash string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.userHistoryURLs = append(f.userHistoryURLs, UserHistoryURL{Hash: hash, CookieID: cookieID})
	return nil
}

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

func (f *File) Ping() (err error) {
	return err
}
