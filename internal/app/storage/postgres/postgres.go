package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
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
	return db.Ping(context.Background())
}

//type Book struct {
//	Isbn   string
//	Title  string
//	Author string
//	Price  float32
//}

func Close() {
	log.Print("DB connection closed.")
	db.Close(context.Background())
}

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
