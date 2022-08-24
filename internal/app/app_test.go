package app

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
				response: `http://localhost:8080/e9db20b246fb7d3ffba1b2182fbcf167`,
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
				response: `http://localhost:8080/e9db20b246fb7d3ffba1b2182fbcf167`,
				postData: `https://yandex.ru`,
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {

			// делаем тестовый http запрос
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http://localhost:8080", bytes.NewBuffer([]byte(tt.want.postData)))

			application := New()
			application.SaveURLHandler(w, r)

			res := w.Result()
			defer res.Body.Close()

			// проверяем код ответа
			require.Equal(t, tt.want.code, res.StatusCode)

			// получаем и проверяем тело запроса
			defer res.Body.Close()
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
			request: "http://localhost:8080/e9db20b246fb7d3ffba1b2182fbcf167",
			want: want{
				code:     307,
				response: `https://yandex.ru`,
			},
		},
		{
			name:    "GET Просмотр несуществующей ссылки",
			request: "http://localhost:8080/0000000000",
			want: want{
				code: 400,
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {

			// делаем тестовый http запрос
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", tt.request, nil)

			app := New()
			app.AddURL("https://yandex.ru")
			app.GetURLHandler(w, r)

			res := w.Result()
			defer res.Body.Close()

			// проверяем код ответа
			//require.Equal(t, tt.want.code, res.StatusCode)

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
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
				url:  `http://localhost:8080/api/shorten`,
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
				url:  `http://localhost:8080/api/shorten`,
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
				url:  `http://localhost:8080/api/shorten`,
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

			// делаем тестовый http запрос
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", tt.request.url, bytes.NewBuffer([]byte(tt.request.data)))

			application := New()
			application.SaveURLJSONHandler(w, r)

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

	type args struct {
		w       http.ResponseWriter
		code    int
		message string
	}

	w := httptest.NewRecorder()

	tests := []struct {
		name string
		args args
	}{
		{
			name: "stringCode",
			args: args{
				w:       w,
				code:    400,
				message: "Only Json format required in request body",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a := &App{}
			a.ShowJSONError(tt.args.w, tt.args.code, tt.args.message)
		})
	}
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
			want: "http://localhost:8080/tescode",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{}
			if got := a.GenerateShortenURL(tt.args.shortenCode); got != tt.want {
				t.Errorf("GenerateShortenURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
