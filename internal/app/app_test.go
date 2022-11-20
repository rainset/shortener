package app

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/rainset/shortener/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var conf = Config{
	ServerAddress:  "localhost:8080",
	ServerBaseURL:  "http://localhost:8080",
	CookieHashKey:  "49a8aca82c132d8d1f430e32be1e6ff3",
	CookieBlockKey: "49a8aca82c132d8d1f430e32be1e6ff2",
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {

	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
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

	client := resty.New()
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	longURL := "http://yandex.ru"

	resp, _ := client.R().
		SetBody(longURL).
		Post(app.Config.ServerBaseURL)

	url1 := resp.String()

	require.Equal(t, http.StatusCreated, resp.StatusCode())

	resp2, _ := client.R().Get(url1)

	require.Equal(t, http.StatusTemporaryRedirect, resp2.StatusCode())
	require.Equal(t, longURL, resp2.Header().Get("Location"))
}

func TestApp_SaveURLJSONHandler(t *testing.T) {

	type request struct {
		url  string
		data string
	}

	type want struct {
		code     int
		response string
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
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readBodyBytes(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("readBodyBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readBodyBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}
