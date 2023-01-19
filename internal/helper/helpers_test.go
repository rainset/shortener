package helper

import (
	"reflect"
	"testing"
)

func BenchmarkGenerateRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateRandom(32)
	}
}

func BenchmarkGenerateToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateToken(8)
	}
}

func TestGenerateRandom(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateRandom(tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRandom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateRandom() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateShortenURL(t *testing.T) {
	type args struct {
		baseURL     string
		shortenCode string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateShortenURL(tt.args.baseURL, tt.args.shortenCode); got != tt.want {
				t.Errorf("GenerateShortenURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateToken(tt.args.length); got != tt.want {
				t.Errorf("GenerateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateUniqueuserID(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateUniqueuserID(); got != tt.want {
				t.Errorf("GenerateUniqueuserID() = %v, want %v", got, tt.want)
			}
		})
	}
}
