package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
)

// ключ подписи
const secretAesKey = "49a8aca82c132d8d1f430e32be1e6ff3"

func GenerateToken(length int) string {
	token := ""
	codeAlphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	codeAlphabet += "abcdefghijklmnopqrstuvwxyz"
	codeAlphabet += "0123456789"

	for i := 0; i < length; i++ {
		token += string(codeAlphabet[cryptoRandSecure(int64(len(codeAlphabet)))])
	}
	return token
}

func cryptoRandSecure(max int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		log.Println(err)
	}
	return nBig.Int64()
}

func GenerateRandom(size int) ([]byte, error) {
	// генерируем случайную последовательность байт
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func EncryptString(data string) (string, error) {
	// подписываемое сообщение
	src := []byte(data)
	secretAesKey := []byte(secretAesKey)
	//fmt.Printf("original: %s\n", src)

	aesblock, err := aes.NewCipher(secretAesKey)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}

	// создаём вектор инициализации
	nonce, err := GenerateRandom(aesgcm.NonceSize())
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}

	dst := aesgcm.Seal(nil, nonce, src, nil) // зашифровываем
	complexHash := fmt.Sprintf("%x.%x", nonce, dst)

	fmt.Println(complexHash)

	return complexHash, nil
}

func DecryptString(complexHash string) (decryptedString string, err error) {

	//// ключ подписи
	secretAesKey := []byte(secretAesKey)

	aesblock, err := aes.NewCipher(secretAesKey)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}

	s := strings.Split(complexHash, ".")

	nonce, _ := hex.DecodeString(s[0])
	origHash, _ := hex.DecodeString(s[1])

	src2, err := aesgcm.Open(nil, nonce, origHash, nil) // расшифровываем
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}

	return string(src2), nil
}
