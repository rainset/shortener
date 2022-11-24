package app

import (
	"fmt"
	"github.com/rainset/shortener/internal/storage/file"
	"github.com/rainset/shortener/internal/storage/postgres"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rainset/shortener/internal/storage/memory"
)

func TestApp_GenerateShortenURL(t *testing.T) {
	hash := "testHash"
	s := memory.New()
	app := New(s, conf)
	shortenURL := app.GenerateShortenURL(hash)
	result := app.Config.ServerBaseURL + "/" + hash
	assert.Equal(t, shortenURL, result)
}

// Пример запуска приложения с бд postgres
func Example_posgres() {

	var postgresDSN = "postgres://root:12345@localhost:5432/shorten"

	conf := Config{
		ServerAddress:  "localhost:8080",
		ServerBaseURL:  "http://localhost",
		CookieHashKey:  "49a8aca82c132d8d1f430e32be1e6ff3",
		CookieBlockKey: "49a8aca82c132d8d1f430e32be1e6ff2",
	}

	s := postgres.New(postgresDSN)
	a := New(s, conf)
	shortenURL := a.GenerateShortenURL("dk8Sv98F")
	fmt.Println(shortenURL)

	// Output:
	// http://localhost:8080/dk8Sv98F
}

// Пример запуска приложения с бд в оперативной памяти
func Example_memory() {

	conf := Config{
		ServerAddress:  "localhost:8080",
		ServerBaseURL:  "http://localhost",
		CookieHashKey:  "49a8aca82c132d8d1f430e32be1e6ff3",
		CookieBlockKey: "49a8aca82c132d8d1f430e32be1e6ff2",
	}

	s := memory.New()
	a := New(s, conf)
	shortenURL := a.GenerateShortenURL("dk8Sv98F")
	fmt.Println(shortenURL)

	// Output:
	// http://localhost:8080/dk8Sv98F
}

// Пример запуска приложения с бд в файловой системе
func Example_files() {

	fileStoragePath := "/storage.txt"

	conf := Config{
		ServerAddress:  "localhost:8080",
		ServerBaseURL:  "http://localhost",
		CookieHashKey:  "49a8aca82c132d8d1f430e32be1e6ff3",
		CookieBlockKey: "49a8aca82c132d8d1f430e32be1e6ff2",
	}

	s := file.New(fileStoragePath)
	a := New(s, conf)
	shortenURL := a.GenerateShortenURL("dk8Sv98F")
	fmt.Println(shortenURL)

	// Output:
	// http://localhost:8080/dk8Sv98F
}

func ExampleApp_GenerateShortenURL() {
	s := memory.New()
	a := New(s, conf)
	shortenURL := a.GenerateShortenURL("dk8Sv98F")
	fmt.Println(shortenURL)

	// Output:
	// http://localhost:8080/dk8Sv98F
}
