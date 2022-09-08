package main

import (
	"compress/gzip"
	"github.com/gorilla/handlers"
	"github.com/rainset/shortener/internal/app"
	"log"
	"net/http"
)

func main() {

	//now := time.Now()
	//sec := now.Unix()
	//
	//fmt.Println(sec)
	//
	//rnd, err := helper.GenerateRandom(32)
	//fmt.Println(rnd, err)
	//
	//strRnd := hex.EncodeToString(rnd)
	//
	//fmt.Println(strRnd)
	//return

	application := app.New()
	application.InitFlags()

	//complexHash, _ := helper.EncryptString("Привет как дела!")
	//
	//fmt.Println(complexHash)
	//
	//helper.DecryptString(complexHash)
	//decryptedString, err := helper.DecryptString(complexHash)
	//
	//fmt.Println(decryptedString, err)

	r := application.NewRouter()
	http.Handle("/", r)

	log.Printf("Listening %s ...", application.Config.ServerAddress)
	log.Fatal(http.ListenAndServe(application.Config.ServerAddress, handlers.CompressHandlerLevel(r, gzip.BestSpeed)))

}
