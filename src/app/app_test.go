package app

import (
	"fmt"
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

func ExampleApp_GenerateShortenURL() {
	s := memory.New()
	a := New(s, conf)
	shortenURL := a.GenerateShortenURL("dk8Sv98F")
	fmt.Println(shortenURL)

	// Output:
	// http://localhost:8080/dk8Sv98F
}
