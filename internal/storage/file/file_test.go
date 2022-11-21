package file

import (
	"github.com/rainset/shortener/internal/storage"
	"os"
	"reflect"
	"sync"
	"testing"
)

func TestFile_AddBatchURL(t *testing.T) {
	type fields struct {
		mutex           sync.Mutex
		fileStoragePath string
		userHistoryURLs []UserHistoryURL
	}
	type args struct {
		in0 []storage.BatchUrls
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult []storage.ResultBatchUrls
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				mutex:           tt.fields.mutex,
				fileStoragePath: tt.fields.fileStoragePath,
				userHistoryURLs: tt.fields.userHistoryURLs,
			}
			gotResult, err := f.AddBatchURL(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddBatchURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("AddBatchURL() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestFile_AddURL(t *testing.T) {
	type fields struct {
		mutex           sync.Mutex
		fileStoragePath string
		userHistoryURLs []UserHistoryURL
	}
	type args struct {
		hash     string
		original string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				mutex:           tt.fields.mutex,
				fileStoragePath: tt.fields.fileStoragePath,
				userHistoryURLs: tt.fields.userHistoryURLs,
			}
			if err := f.AddURL(tt.args.hash, tt.args.original); (err != nil) != tt.wantErr {
				t.Errorf("AddURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFile_AddUserHistoryURL(t *testing.T) {
	type fields struct {
		mutex           sync.Mutex
		fileStoragePath string
		userHistoryURLs []UserHistoryURL
	}
	type args struct {
		cookieID string
		hash     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				mutex:           tt.fields.mutex,
				fileStoragePath: tt.fields.fileStoragePath,
				userHistoryURLs: tt.fields.userHistoryURLs,
			}
			if err := f.AddUserHistoryURL(tt.args.cookieID, tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("AddUserHistoryURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFile_DeleteBatchURL(t *testing.T) {
	type fields struct {
		mutex           sync.Mutex
		fileStoragePath string
		userHistoryURLs []UserHistoryURL
	}
	type args struct {
		in0 []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				mutex:           tt.fields.mutex,
				fileStoragePath: tt.fields.fileStoragePath,
				userHistoryURLs: tt.fields.userHistoryURLs,
			}
			if err := f.DeleteBatchURL(tt.args.in0); (err != nil) != tt.wantErr {
				t.Errorf("DeleteBatchURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFile_DeleteUserBatchURL(t *testing.T) {
	type fields struct {
		mutex           sync.Mutex
		fileStoragePath string
		userHistoryURLs []UserHistoryURL
	}
	type args struct {
		in0 string
		in1 []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				mutex:           tt.fields.mutex,
				fileStoragePath: tt.fields.fileStoragePath,
				userHistoryURLs: tt.fields.userHistoryURLs,
			}
			if err := f.DeleteUserBatchURL(tt.args.in0, tt.args.in1); (err != nil) != tt.wantErr {
				t.Errorf("DeleteUserBatchURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFile_GetByOriginalURL(t *testing.T) {
	type fields struct {
		mutex           sync.Mutex
		fileStoragePath string
		userHistoryURLs []UserHistoryURL
	}
	type args struct {
		original string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantHash string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				mutex:           tt.fields.mutex,
				fileStoragePath: tt.fields.fileStoragePath,
				userHistoryURLs: tt.fields.userHistoryURLs,
			}
			gotHash, err := f.GetByOriginalURL(tt.args.original)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByOriginalURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHash != tt.wantHash {
				t.Errorf("GetByOriginalURL() gotHash = %v, want %v", gotHash, tt.wantHash)
			}
		})
	}
}

func TestFile_GetListUserHistoryURL(t *testing.T) {
	type fields struct {
		mutex           sync.Mutex
		fileStoragePath string
		userHistoryURLs []UserHistoryURL
	}
	type args struct {
		cookieID string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult []storage.ResultHistoryURL
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				mutex:           tt.fields.mutex,
				fileStoragePath: tt.fields.fileStoragePath,
				userHistoryURLs: tt.fields.userHistoryURLs,
			}
			gotResult, err := f.GetListUserHistoryURL(tt.args.cookieID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListUserHistoryURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("GetListUserHistoryURL() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestFile_GetURL(t *testing.T) {
	type fields struct {
		mutex           sync.Mutex
		fileStoragePath string
		userHistoryURLs []UserHistoryURL
	}
	type args struct {
		hash string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantResultURL storage.ResultURL
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				mutex:           tt.fields.mutex,
				fileStoragePath: tt.fields.fileStoragePath,
				userHistoryURLs: tt.fields.userHistoryURLs,
			}
			gotResultURL, err := f.GetURL(tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResultURL, tt.wantResultURL) {
				t.Errorf("GetURL() gotResultURL = %v, want %v", gotResultURL, tt.wantResultURL)
			}
		})
	}
}

func TestFile_Ping(t *testing.T) {
	type fields struct {
		mutex           sync.Mutex
		fileStoragePath string
		userHistoryURLs []UserHistoryURL
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				mutex:           tt.fields.mutex,
				fileStoragePath: tt.fields.fileStoragePath,
				userHistoryURLs: tt.fields.userHistoryURLs,
			}
			if err := f.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		fileStoragePath string
	}
	tests := []struct {
		name string
		args args
		want *File
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.fileStoragePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConsumer(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    *consumer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConsumer(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConsumer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConsumer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewProducer(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    *producer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProducer(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProducer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProducer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_consumer_Close(t *testing.T) {
	type fields struct {
		file    *os.File
		decoder *json.Decoder
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &consumer{
				file:    tt.fields.file,
				decoder: tt.fields.decoder,
			}
			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_consumer_ReadURL(t *testing.T) {
	type fields struct {
		file    *os.File
		decoder *json.Decoder
	}
	tests := []struct {
		name    string
		fields  fields
		want    *DataURL
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &consumer{
				file:    tt.fields.file,
				decoder: tt.fields.decoder,
			}
			got, err := c.ReadURL()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_consumer_RestoreStorage(t *testing.T) {
	type fields struct {
		file    *os.File
		decoder *json.Decoder
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult []ResultURL
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &consumer{
				file:    tt.fields.file,
				decoder: tt.fields.decoder,
			}
			gotResult, err := c.RestoreStorage()
			if (err != nil) != tt.wantErr {
				t.Errorf("RestoreStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("RestoreStorage() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_producer_Close(t *testing.T) {
	type fields struct {
		file    *os.File
		encoder *json.Encoder
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &producer{
				file:    tt.fields.file,
				encoder: tt.fields.encoder,
			}
			if err := p.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_producer_WriteURL(t *testing.T) {
	type fields struct {
		file    *os.File
		encoder *json.Encoder
	}
	type args struct {
		url *DataURL
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &producer{
				file:    tt.fields.file,
				encoder: tt.fields.encoder,
			}
			if err := p.WriteURL(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("WriteURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
