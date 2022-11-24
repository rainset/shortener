// Пакет helper вспомогательные функции для проекта
package helper

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"
	"time"
)

// GenerateToken генерирует случайную строку длинной length
func GenerateToken(length int) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return strings.ToLower(base32.StdEncoding.EncodeToString(randomBytes)[:length])
}

// GenerateRandom генерирует случайную последовательность байт длинной size
func GenerateRandom(size int) ([]byte, error) {
	// генерируем случайную последовательность байт
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateUniqueuserID генерирует уникальный идентификатор для пользователя
func GenerateUniqueuserID() string {
	now := time.Now()
	sec := now.Unix()
	rnd, _ := GenerateRandom(32)
	return fmt.Sprintf("user.%d.%x", sec, rnd)
}
