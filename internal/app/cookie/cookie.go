package cookie

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/rainset/shortener/internal/app/helper"
	"net/http"
	"time"
)

// Hash keys should be at least 32 bytes long
var hashKey = []byte("49a8aca82c132d8d1f430e32be1e6ff3")

// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
// Shorter keys may weaken the encryption used.
var blockKey = []byte("49a8aca82c132d8d1f430e32be1e6ff2")
var s = securecookie.New(hashKey, blockKey)

func GenerateUniqueuserID() string {
	now := time.Now()
	sec := now.Unix()
	rnd, _ := helper.GenerateRandom(32)
	return fmt.Sprintf("user.%d.%x", sec, rnd)
}

func Set(w http.ResponseWriter, r *http.Request, cookieName string, cookieValue string) {

	value := map[string]string{
		cookieName: cookieValue,
	}

	if encoded, err := s.Encode(cookieName, &value); err == nil {
		cookie := &http.Cookie{
			Name:     cookieName,
			Value:    encoded,
			Path:     "/",
			Secure:   false,
			HttpOnly: false,
		}
		http.SetCookie(w, cookie)
	}
}

func Get(w http.ResponseWriter, r *http.Request, cookieName string) (value string, err error) {
	if cookie, err := r.Cookie(cookieName); err == nil {
		value := make(map[string]string)
		if err = s.Decode(cookieName, cookie.Value, &value); err == nil {
			return value[cookieName], nil
		}
	}
	return "", err
}
