package app

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {

	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestApp_SaveURLHandler(t *testing.T) {
	// определяем структуру теста
	type want struct {
		code     int
		response string
		postData string
	}

	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		name string
		want want
	}{
		// определяем все тесты
		{
			name: "POST Добавление ссылки",
			want: want{
				code:     201,
				response: `%s/e9db20b246fb7d3ffba1b2182fbcf167`,
				postData: `https://yandex.ru`,
			},
		},
		{
			name: "POST Пустой запрос",
			want: want{
				code:     400,
				response: ``,
				postData: ``,
			},
		},
		{
			name: "POST Добавление ссылки",
			want: want{
				code:     201,
				response: `%s/e9db20b246fb7d3ffba1b2182fbcf167`,
				postData: `https://yandex.ru`,
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {

			app := New()

			// делаем тестовый http запрос
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", app.Config.ServerBaseURL, bytes.NewBuffer([]byte(tt.want.postData)))

			app.SaveURLHandler(w, r)

			res := w.Result()
			defer res.Body.Close()

			// проверяем код ответа
			require.Equal(t, tt.want.code, res.StatusCode)

			// получаем и проверяем тело запроса
			_, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestApp_GetURLHandler(t *testing.T) {
	// определяем структуру теста
	type want struct {
		code     int
		response string
	}

	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		name    string
		request string
		want    want
	}{
		// определяем все тесты
		{
			name:    "GET Просмотр ссылки",
			request: "/e9db20b246fb7d3ffba1b2182fbcf167",
			want: want{
				code: 200,
			},
		},
		{
			name:    "GET Просмотр несуществующей ссылки",
			request: "/0000000000qq",
			want: want{
				code:     400,
				response: "Bad Url\n",
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			app := New()
			_, err := app.AddURL("https://yandex.ru")
			if err != nil {
				t.Error(err)
			}
			r := app.NewRouter()
			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, body := testRequest(t, ts, "GET", tt.request)

			if tt.want.code != 0 {
				require.Equal(t, tt.want.code, resp.StatusCode)
			}
			if tt.want.response != "" {
				require.Equal(t, tt.want.response, body)
			}
		})
	}
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
				response: `{"result":"http://localhost:8080/e9db20b246fb7d3ffba1b2182fbcf167"}`,
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
			app := New()

			// делаем тестовый http запрос
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", fmt.Sprintf(tt.request.url, app.Config.ServerBaseURL), bytes.NewBuffer([]byte(tt.request.data)))

			app.SaveURLJSONHandler(w, r)

			result := w.Result()
			result.Body.Close()

			// проверяем код ответа
			require.Equal(t, tt.want.code, result.StatusCode)

			data, err := io.ReadAll(result.Body)
			if err != nil {
				t.Fatal(err)
			}
			require.Equal(t, tt.want.response, string(data))

		})
	}
}

func TestApp_ShowJSONError(t *testing.T) {
	w := httptest.NewRecorder()
	a := &App{}
	a.ShowJSONError(w, 400, "Only Json format required in request body")
}

func TestApp_GenerateShortenURL(t *testing.T) {
	type args struct {
		shortenCode string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Generate shorten link",
			args: args{
				shortenCode: "tescode",
			},
			want: "%s/tescode",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := New()
			need := fmt.Sprintf(tt.want, app.Config.ServerBaseURL)
			if got := app.GenerateShortenURL(tt.args.shortenCode); got != need {
				t.Errorf("GenerateShortenURL() = %v, want %v", got, need)
			}
		})
	}
}
