package cookie

import (
	"github.com/gorilla/securecookie"
	"net/http"
)

// Hash keys should be at least 32 bytes long
// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
// Shorter keys may weaken the encryption used.

type SCookie struct {
	HasKey   string
	BlockKey string
	s        *securecookie.SecureCookie
}

func New(hashKey, blockKey string) *SCookie {
	s := securecookie.New([]byte(hashKey), []byte(blockKey))
	return &SCookie{
		HasKey:   hashKey,
		BlockKey: blockKey,
		s:        s,
	}
}

func (c *SCookie) Set(w http.ResponseWriter, _ *http.Request, cookieName string, cookieValue string) {

	value := map[string]string{
		cookieName: cookieValue,
	}

	if encoded, err := c.s.Encode(cookieName, &value); err == nil {
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

func (c *SCookie) Get(_ http.ResponseWriter, r *http.Request, cookieName string) (value string, err error) {
	if cookie, err := r.Cookie(cookieName); err == nil {
		value := make(map[string]string)
		if err = c.s.Decode(cookieName, cookie.Value, &value); err == nil {
			return value[cookieName], nil
		}
	}
	return value, err
}
