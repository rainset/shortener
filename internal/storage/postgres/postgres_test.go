package postgres

import (
	"github.com/rainset/shortener/internal/storage"
	"reflect"
	"testing"
)

func TestCreateTables(t *testing.T) {
	type args struct {
		db *pgx.Conn
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateTables(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("CreateTables() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabase_AddBatchURL(t *testing.T) {
	type fields struct {
		pgx *pgx.Conn
	}
	type args struct {
		urls []storage.BatchUrls
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
			d := &Database{
				pgx: tt.fields.pgx,
			}
			gotResult, err := d.AddBatchURL(tt.args.urls)
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

func TestDatabase_AddURL(t *testing.T) {
	type fields struct {
		pgx *pgx.Conn
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
			d := &Database{
				pgx: tt.fields.pgx,
			}
			if err := d.AddURL(tt.args.hash, tt.args.original); (err != nil) != tt.wantErr {
				t.Errorf("AddURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabase_AddUserHistoryURL(t *testing.T) {
	type fields struct {
		pgx *pgx.Conn
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
			d := &Database{
				pgx: tt.fields.pgx,
			}
			if err := d.AddUserHistoryURL(tt.args.cookieID, tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("AddUserHistoryURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabase_Close(t *testing.T) {
	type fields struct {
		pgx *pgx.Conn
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				pgx: tt.fields.pgx,
			}
			d.Close()
		})
	}
}

func TestDatabase_DeleteBatchURL(t *testing.T) {
	type fields struct {
		pgx *pgx.Conn
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
			d := &Database{
				pgx: tt.fields.pgx,
			}
			if err := d.DeleteBatchURL(tt.args.hashes); (err != nil) != tt.wantErr {
				t.Errorf("DeleteBatchURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabase_DeleteUserBatchURL(t *testing.T) {
	type fields struct {
		pgx *pgx.Conn
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
			d := &Database{
				pgx: tt.fields.pgx,
			}
			if err := d.DeleteUserBatchURL(tt.args.cookieID, tt.args.hashes); (err != nil) != tt.wantErr {
				t.Errorf("DeleteUserBatchURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabase_GetByOriginalURL(t *testing.T) {
	type fields struct {
		pgx *pgx.Conn
	}
	type args struct {
		originalURL string
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
			d := &Database{
				pgx: tt.fields.pgx,
			}
			gotHash, err := d.GetByOriginalURL(tt.args.originalURL)
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

func TestDatabase_GetListUserHistoryURL(t *testing.T) {
	type fields struct {
		pgx *pgx.Conn
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
			d := &Database{
				pgx: tt.fields.pgx,
			}
			gotResult, err := d.GetListUserHistoryURL(tt.args.cookieID)
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

func TestDatabase_GetURL(t *testing.T) {
	type fields struct {
		pgx *pgx.Conn
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
			d := &Database{
				pgx: tt.fields.pgx,
			}
			gotResultURL, err := d.GetURL(tt.args.hash)
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

func TestDatabase_Ping(t *testing.T) {
	type fields struct {
		pgx *pgx.Conn
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
			d := &Database{
				pgx: tt.fields.pgx,
			}
			if err := d.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		dataSourceName string
	}
	tests := []struct {
		name string
		args args
		want *Database
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.dataSourceName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
