package main

import (
	"github.com/rainset/shortener/internal/app"
	"log"
	"net/http"
)

func main() {
	application := app.New()
	application.InitFlags()
	r := application.NewRouter()

	http.Handle("/", r)
	log.Printf("Listening %s ...", application.Config.ServerAddress)
	log.Fatal(http.ListenAndServe(application.Config.ServerAddress, r))
}

//type gzipWriter struct {
//	http.ResponseWriter
//	Writer io.Writer
//}
//
//func (w gzipWriter) Write(b []byte) (int, error) {
//	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
//	return w.Writer.Write(b)
//}

//func gzipHandle(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// проверяем, что клиент поддерживает gzip-сжатие
//		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
//			// если gzip не поддерживается, передаём управление
//			// дальше без изменений
//			next.ServeHTTP(w, r)
//			return
//		}
//
//		// создаём gzip.Writer поверх текущего w
//		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
//		if err != nil {
//			io.WriteString(w, err.Error())
//			return
//		}
//		defer gz.Close()
//
//		w.Header().Set("Content-Encoding", "gzip")
//		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
//		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
//	})
//}
