package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/rainset/shortener/internal/helper"
	"github.com/rainset/shortener/internal/storage"
	"io/ioutil"
	"log"
)

type Database struct {
	storage.InterfaceStorage
	pgx *pgx.Conn
}
type ResultURL struct {
	ID       int
	Hash     string
	Original string
}
type UserHistoryURL struct {
	ID       int
	CookieID string
	Hash     string
}

func Init(dataSourceName string) *Database {
	db, err := pgx.Connect(context.Background(), dataSourceName)

	if db == nil && err == nil {
		err = errors.New("connection problems")
	}

	if err == nil {
		log.Print("DB connection initialized...")
		err = CreateTables(db)
		if err != nil {
			log.Println(err)
		}
	}

	return &Database{
		pgx: db,
	}
}

func CreateTables(db *pgx.Conn) error {

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

func (d *Database) Ping() error {

	if d.pgx == nil {
		return errors.New("connection not initialized")
	}

	err := d.pgx.Ping(context.Background())
	if err != nil {
		log.Printf("ping error: %s", err)
	}
	return err
}

func (d *Database) Close() {
	log.Print("DB connection closed.")
	err := d.pgx.Close(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}

func (d *Database) GetURL(hash string) (string, error) {
	var item ResultURL
	q := "SELECT * FROM urls WHERE hash = $1"
	err := d.pgx.QueryRow(context.Background(), q, hash).Scan(&item.ID, &item.Hash, &item.Original)

	if err != nil {
		return item.Original, err
	}
	return item.Original, nil
}

func (d *Database) GetByOriginalURL(originalURL string) (hash string, err error) {
	var item ResultURL
	q := "SELECT * FROM urls WHERE original = $1"
	err = d.pgx.QueryRow(context.Background(), q, originalURL).Scan(&item.ID, &item.Hash, &item.Original)

	if err != nil {
		return "", err
	}
	return item.Hash, nil
}

func (d *Database) AddURL(hash, original string) error {
	q := "INSERT INTO urls (hash, original) VALUES ($1, $2)"
	_, err := d.pgx.Exec(context.Background(), q, hash, original)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) AddUserHistoryURL(cookieID, hash string) error {
	q := "INSERT INTO user_history_urls (cookie_id, hash) VALUES ($1, $2)"
	_, err := d.pgx.Exec(context.Background(), q, cookieID, hash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				err = nil
			}
		}
		return err
	}
	return nil
}

func (d *Database) GetListUserHistoryURL(cookieID string) (result []storage.ResultHistoryURL, err error) {

	q := "SELECT DISTINCT uh.hash, uh.id, uh.cookie_id, u.original FROM user_history_urls uh INNER JOIN urls u ON u.hash = uh.hash WHERE uh.cookie_id =$1"
	rows, err := d.pgx.Query(context.Background(), q, cookieID)

	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		rowArray := storage.ResultHistoryURL{}
		err := rows.Scan(&rowArray.Hash, &rowArray.ID, &rowArray.CookieID, &rowArray.Original)
		if err != nil {
			return result, err
		}
		result = append(result, storage.ResultHistoryURL{ID: rowArray.ID, Hash: rowArray.Hash, CookieID: rowArray.CookieID, Original: rowArray.Original})
	}
	return result, err
}

func (d *Database) AddBatchURL(urls []storage.BatchUrls) (result []storage.ResultBatchUrls, err error) {

	tx, err := d.pgx.Begin(context.Background())
	if err != nil {
		return result, err
	}

	q := "INSERT INTO urls (hash, original) VALUES ($1, $2)"
	_, err = tx.Prepare(context.Background(), "batch_insert", q)
	if err != nil {
		return result, err
	}

	var hash string
	for _, v := range urls {
		hash = helper.GenerateToken(8)
		_, err = tx.Exec(context.Background(), "batch_insert", hash, v.OriginalURL)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.UniqueViolation {
					err = nil
					//hash, errItem := d.GetByOriginalURL(v.OriginalURL)
					//if errItem != nil {
					//	return result, errItem
					//}
				}
			}
		}
		if err == nil {
			result = append(result, storage.ResultBatchUrls{CorrelationID: v.CorrelationID, Hash: hash})
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
