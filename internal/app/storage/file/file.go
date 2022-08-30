package file

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type DataURL struct {
	Hash    string `json:"hash"`
	LongUrl string `json:"long_url"`
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
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

type consumer struct {
	file    *os.File
	decoder *json.Decoder
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

func (c *consumer) RestoreStorage() (map[string]string, error) {
	scanner := bufio.NewScanner(c.file)
	readedURLs := make(map[string]string)
	for scanner.Scan() {
		dataURL := &DataURL{}
		err := json.Unmarshal([]byte(scanner.Text()), dataURL)
		if err != nil {
			log.Fatal(err)
		}
		readedURLs[dataURL.Hash] = dataURL.LongUrl
	}

	return readedURLs, nil
}
