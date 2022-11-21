package app

import (
	"github.com/rainset/shortener/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApp_GenerateShortenURL(t *testing.T) {
	hash := "testHash"
	s := memory.New()
	app := New(s, conf)
	shortenURL := app.GenerateShortenURL(hash)
	result := app.Config.ServerBaseURL + "/" + hash
	assert.Equal(t, shortenURL, result)
}
