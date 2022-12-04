package main

import (
	"encoding/base32"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"math/rand"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"time"
)

func main() {

	// создаем cookie jar для сохранения cookies между запросами
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic("Неожиданная ошибка при создании Cookie Jar")
	}

	client := resty.New().SetCookieJar(jar)
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	for i := 0; i < 500; i++ {

		randS := fmt.Sprintf("%s-%d", getToken(32), time.Now().Unix())

		resp0, err := client.R().Get("http://localhost:8080/ping")
		if err != nil {
			panic(err)
		}
		log.Println("resp0", resp0.StatusCode())

		resp1, err := client.R().SetBody(fmt.Sprintf(`http://google.com/?%d/%s/`, i, randS)).Post("http://localhost:8080/")
		if err != nil {
			panic(err)
		}
		log.Println("resp1", resp1.StatusCode(), resp1.String())

		resp2, _ := client.R().Get(resp1.String())

		log.Println("resp2", resp2.StatusCode())

		type testBatchData struct {
			CorrelationID string `json:"correlation_id"`
			OriginalURL   string `json:"original_url"`
		}
		var arrBatch []testBatchData
		for c := 1; c < 10; c++ {
			arrBatch = append(arrBatch, testBatchData{CorrelationID: strconv.Itoa(c), OriginalURL: fmt.Sprintf("http://ya.ru/%d/%s", c, randS)})
		}

		resp3, _ := client.R().SetBody(arrBatch).Post("http://localhost:8080/api/shorten/batch")
		log.Println("resp3", resp3.StatusCode(), resp3)

		//resp4, err := client.R().Get("http://localhost:8080/api/user/urls")
		//if err != nil {
		//	panic(err)
		//}
		//log.Println("resp4", resp4.StatusCode(), resp4)
	}

}

func getToken(length int) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return strings.ToLower(base32.StdEncoding.EncodeToString(randomBytes)[:length])
}
