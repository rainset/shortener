package app

import (
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
