package memory

import (
	"github.com/rainset/shortener/internal/storage"
	"reflect"
	"sync"
	"testing"
)

func TestMemory_AddBatchURL(t *testing.T) {
	type fields struct {
		mutex           sync.RWMutex
		urls            map[string]storage.ResultURL
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
			m := &Memory{
				mutex:           tt.fields.mutex,
				urls:            tt.fields.urls,
				userHistoryURLs: tt.fields.userHistoryURLs,
			}
			gotResult, err := m.AddBatchURL(tt.args.in0)
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

func TestMemory_AddURL(t *testing.T) {
	type fields struct {
		mutex           sync.RWMutex
		urls            map[string]storage.ResultURL
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
			m := &Memory{
				mutex:           tt.fields.mutex,
				urls:            tt.fields.urls,
				userHistoryURLs: tt.fields.userHistoryURLs,
			}
			if err := m.AddURL(tt.args.hash, tt.args.original); (err != nil) != tt.wantErr {
				t.Errorf("AddURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_AddUserHistoryURL(t *testing.T) {
	type fields struct {
		mutex           sync.RWMutex
		urls            map[string]storage.ResultURL
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
			m := &Memory{
				mutex:           tt.fields.mutex,
				urls:            tt.fields.urls,
				userHistoryURLs: tt.fields.userHistoryURLs,
			}
			if err := m.AddUserHistoryURL(tt.args.cookieID, tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("AddUserHistoryURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_DeleteBatchURL(t *testing.T) {
	type fields struct {
		mutex           sync.RWMutex
		urls            map[string]storage.ResultURL
		userHistoryURLs []UserHistoryURL
	}
	type args struct {
		hashes []string
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
			m := &Memory{
				mutex:           tt.fields.mutex,
				urls:            tt.fields.urls,
				userHistoryURLs: tt.fields.userHistoryURLs,
			}
			if err := m.DeleteBatchURL(tt.args.hashes); (err != nil) != tt.wantErr {
				t.Errorf("DeleteBatchURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_DeleteUserBatchURL(t *testing.T) {
	type fields struct {
		mutex           sync.RWMutex
		urls            map[string]storage.ResultURL
		userHistoryURLs []UserHistoryURL
	}
	type args struct {
		cookieID string
		hashes   []string
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
			m := &Memory{
				mutex:           tt.fields.mutex,
				urls:            tt.fields.urls,
				userHistoryURLs: tt.fields.userHistoryURLs,
			}
			if err := m.DeleteUserBatchURL(tt.args.cookieID, tt.args.hashes); (err != nil) != tt.wantErr {
				t.Errorf("DeleteUserBatchURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_GetByOriginalURL(t *testing.T) {
	type fields struct {
		mutex           sync.RWMutex
		urls            map[string]storage.ResultURL
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
			m := &Memory{
				mutex:           tt.fields.mutex,
				urls:            tt.fields.urls,
				userHistoryURLs: tt.fields.userHistoryURLs,
			}
			gotHash, err := m.GetByOriginalURL(tt.args.original)
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

func TestMemory_GetListUserHistoryURL(t *testing.T) {
	type fields struct {
		mutex           sync.RWMutex
		urls            map[string]storage.ResultURL
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
			m := &Memory{
				mutex:           tt.fields.mutex,
				urls:            tt.fields.urls,
				userHistoryURLs: tt.fields.userHistoryURLs,
			}
			gotResult, err := m.GetListUserHistoryURL(tt.args.cookieID)
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

func TestMemory_GetURL(t *testing.T) {
	type fields struct {
		mutex           sync.RWMutex
		urls            map[string]storage.ResultURL
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
			m := &Memory{
				mutex:           tt.fields.mutex,
				urls:            tt.fields.urls,
				userHistoryURLs: tt.fields.userHistoryURLs,
			}
			gotResultURL, err := m.GetURL(tt.args.hash)
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

func TestMemory_Ping(t *testing.T) {
	type fields struct {
		mutex           sync.RWMutex
		urls            map[string]storage.ResultURL
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
			m := &Memory{
				mutex:           tt.fields.mutex,
				urls:            tt.fields.urls,
				userHistoryURLs: tt.fields.userHistoryURLs,
			}
			if err := m.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Memory
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
