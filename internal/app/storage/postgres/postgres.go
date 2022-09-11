package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"io/ioutil"
	"log"
)

// This time the global variable is unexported.
var db *pgx.Conn

// InitDB sets up setting up the connection pool global variable.
func InitDB(dataSourceName string) (err error) {
	db, err = pgx.Connect(context.Background(), dataSourceName)

	if db == nil && err == nil {
		err = errors.New("connection problems")
	}

	if err != nil {
		log.Println(err)
		return err
	}
	log.Print("DB connection initialized...")
	return err
}

func Ping() {
	err := db.Ping(context.Background())
	if err != nil {
		log.Printf("ping error: %s", err)
	}
}

func Close() {
	log.Print("DB connection closed.")
	db.Close(context.Background())
}

func CreateTables() {
	c, ioErr := ioutil.ReadFile("migrations/tables.sql")
	if ioErr != nil {
		log.Println("read file tables: ", ioErr)
	}
	q := string(c)
	_, err := db.Exec(context.Background(), q)
	if err != nil {
		log.Println("create tables: ", err)
	}
	log.Println("tables created")
}

func AddURL(hash, originalURL string) error {
	q := "INSERT INTO urls (hash, original) VALUES ($1, $2)"
	_, err := db.Exec(context.Background(), q, hash, originalURL)
	if err != nil {
		return err
	}
	return nil
}

//type Book struct {
//	Isbn   string
//	Title  string
//	Author string
//	Price  float32
//}

//func AllBooks() ([]Book, error) {
//	// This now uses the unexported global variable.
//	rows, err := db.Query("SELECT * FROM books")
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var bks []Book
//
//	for rows.Next() {
//		var bk Book
//
//		err := rows.Scan(&bk.Isbn, &bk.Title, &bk.Author, &bk.Price)
//		if err != nil {
//			return nil, err
//		}
//
//		bks = append(bks, bk)
//	}
//	if err = rows.Err(); err != nil {
//		return nil, err
//	}
//
//	return bks, nil
//}
