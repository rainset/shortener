package helper

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"
)

func GenerateToken(length int) string {
	//token := ""
	//codeAlphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	//codeAlphabet += "abcdefghijklmnopqrstuvwxyz"
	//codeAlphabet += "0123456789"
	//
	//for i := 0; i < length; i++ {
	//	token += string(codeAlphabet[cryptoRandSecure(int64(len(codeAlphabet)))])
	//}
	//return token
	//
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return strings.ToLower(base32.StdEncoding.EncodeToString(randomBytes)[:length])
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

func GenerateUniqueuserID() string {
	now := time.Now()
	sec := now.Unix()
	rnd, _ := GenerateRandom(32)
	return fmt.Sprintf("user.%d.%x", sec, rnd)
}
