package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/rainset/shortener/internal/app/helper"
	"io/ioutil"
	"log"
)

// This time the global variable is unexported.
var db *pgx.Conn

type BatchUrls struct {
	CorrelationID string
	OriginalURL   string
}

type ResultBatchUrls struct {
	CorrelationID string
	Hash          string
}

type ResultURL struct {
	ID       int
	Hash     string
	Original string
}

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

func CreateTables() error {

	//if db == nil {
	//	return errors.New("create tables, err connection")
	//}

	c, ioErr := ioutil.ReadFile("migrations/tables.sql")
	if ioErr != nil {
		log.Println("read file tables: ", ioErr)
		return ioErr
	}
	q := string(c)
	_, err := db.Exec(context.Background(), q)
	if err != nil {
		log.Println("create tables: ", err)
		return err
	}
	log.Println("tables created")
	return nil
}

func GetURL(hash string) (ResultURL, error) {
	var item ResultURL
	q := "SELECT * FROM urls WHERE hash = $1"
	err := db.QueryRow(context.Background(), q, hash).Scan(&item.ID, &item.Hash, &item.Original)

	if err != nil {
		return item, err
	}
	return item, nil
}

func GetByOriginalURL(originalURL string) (ResultURL, error) {
	var item ResultURL
	q := "SELECT * FROM urls WHERE original = $1"
	err := db.QueryRow(context.Background(), q, originalURL).Scan(&item.ID, &item.Hash, &item.Original)

	if err != nil {
		return item, err
	}
	return item, nil
}

func AddURL(hash, originalURL string) error {
	q := "INSERT INTO urls (hash, original) VALUES ($1, $2)"
	_, err := db.Exec(context.Background(), q, hash, originalURL)
	if err != nil {
		return err
	}
	return nil
}

func AddBatchURL(urls *[]BatchUrls) ([]ResultBatchUrls, error) {

	var err error
	var result []ResultBatchUrls

	tx, err := db.Begin(context.Background())
	if err != nil {
		return result, err
	}

	q := "INSERT INTO urls (hash, original) VALUES ($1, $2)"
	_, err = tx.Prepare(context.Background(), "batch_insert", q)
	if err != nil {
		return result, err
	}

	for _, v := range *urls {
		hash := helper.GenerateToken(8)
		_, err = tx.Exec(context.Background(), "batch_insert", hash, v.OriginalURL)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.UniqueViolation {
					err = nil
					item, errItem := GetByOriginalURL(v.OriginalURL)
					if errItem != nil {
						return result, errItem
					}
					hash = item.Hash
				}
			}
		}
		if err == nil {
			result = append(result, ResultBatchUrls{CorrelationID: v.CorrelationID, Hash: hash})
		}
	}

	//if err != nil {
	//	_ = tx.Rollback(context.Background())
	//	return result, err
	//}
	err = tx.Commit(context.Background())
	if err != nil {
		return result, err
	}
	return result, nil
}
