package main

import (
	"fmt"
	app "github.com/rainset/shortener/internal/app"
	"log"
	"net/http"
)

const addr = "localhost:8080"

func main() {

	application := app.New()

	application.AddUrl("link1", "http://yandex.ru")

	fmt.Println(application.GetUrl("link1"))

	http.HandleFunc("/", application.RouteHandler)
	log.Print("Listening...")
	log.Fatal(http.ListenAndServe(addr, nil))
}
