package app

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rainset/shortener/internal/storage/memory"
)

var conf = Config{
	ServerAddress:  "localhost:8080",
	ServerBaseURL:  "http://localhost:8080",
	CookieHashKey:  "49a8aca82c132d8d1f430e32be1e6ff3",
	CookieBlockKey: "49a8aca82c132d8d1f430e32be1e6ff2",
}

func TestApp_SaveURLHandler(t *testing.T) {

	t.Run("POST add url", func(t *testing.T) {
		s := memory.New()
		app := New(s, conf)
		r := app.NewRouter()
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", app.Config.ServerBaseURL+"/", bytes.NewBuffer([]byte("http://yandex.ru")))
		r.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		// проверяем код ответа
		require.Equal(t, 201, w.Code)
		//assert.Equal(t, "pong", w.Body.String())
	})

	t.Run("POST empty body", func(t *testing.T) {
		s := memory.New()
		app := New(s, conf)
		r := app.NewRouter()
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", app.Config.ServerBaseURL+"/", bytes.NewBuffer([]byte("")))
		r.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		// проверяем код ответа
		require.Equal(t, 400, w.Code)
		//assert.Equal(t, "pong", w.Body.String())
	})

}

func TestApp_GetURLHandler(t *testing.T) {

	s := memory.New()
	app := New(s, conf)

	_ = app.s.AddURL("testhash", "http://test.com")

	router := app.NewRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/testhash", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)

}

func TestApp_SaveURLJSONHandler(t *testing.T) {

	type request struct {
		url  string
		data string
	}

	type want struct {
		response string
		code     int
	}

	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "POST - Добавление ссылки",
			request: request{
				url:  `%s/api/shorten`,
				data: `{"url":"https://yandex.ru"}`,
			},
			want: want{
				code:     201,
				response: `{"result":"http://localhost:8080/SZfLgeBS"}`,
			},
		},
		{
			name: "POST - Пустой запрос",
			request: request{
				url:  `%s/api/shorten`,
				data: ``,
			},
			want: want{
				code:     400,
				response: `{"code":400,"message":"Only Json format required in request body"}`,
			},
		},
		{
			name: "POST - неправильный Json формат запроса",
			request: request{
				url:  `%s/api/shorten`,
				data: `testdata`,
			},
			want: want{
				code:     400,
				response: `{"code":400,"message":"Only Json format required in request body"}`,
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			s := memory.New()
			app := New(s, conf)

			// делаем тестовый http запрос
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf(tt.request.url, app.Config.ServerBaseURL), bytes.NewBuffer([]byte(tt.request.data)))
			r := app.NewRouter()
			r.ServeHTTP(w, req)

			result := w.Result()
			result.Body.Close()

			// проверяем код ответа
			require.Equal(t, tt.want.code, result.StatusCode)

			_, err := io.ReadAll(result.Body)
			if err != nil {
				t.Fatal(err)
			}
			//require.Equal(t, tt.want.response, string(data))

		})
	}
}

func Test_readBodyBytes(t *testing.T) {

}

func TestApp_DeleteUserBatchURLHandler(t *testing.T) {
	s := memory.New()
	app := New(s, conf)

	router := app.NewRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/api/user/urls", bytes.NewBuffer([]byte(`["hashtest"]`)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
}

func TestApp_SaveURLBatchJSONHandler(t *testing.T) {
	s := memory.New()
	app := New(s, conf)

	router := app.NewRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/shorten/batch", bytes.NewBuffer([]byte(`[{"correlation_id":"222","original_url":"http://example.com"}]`)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestApp_UserURLListHandler(t *testing.T) {
	s := memory.New()
	app := New(s, conf)

	router := app.NewRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/user/urls", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestApp_NewRouter(t *testing.T) {

}

func TestApp_PingHandler(t *testing.T) {
	s := memory.New()
	app := New(s, conf)

	router := app.NewRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

}
