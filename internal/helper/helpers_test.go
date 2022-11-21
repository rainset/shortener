package helper

import (
	"testing"
)

func TestGenerateRandom(t *testing.T) {
	_, err := GenerateRandom(10)
	if err != nil {
		t.Errorf("GenerateRandom() error = %v, wantErr %v", err, true)
		return
	}
}

func TestGenerateToken(t *testing.T) {
	_ = GenerateToken(10)
}

func TestGenerateUniqueuserID(t *testing.T) {

}

func Test_cryptoRandSecure(t *testing.T) {

}
