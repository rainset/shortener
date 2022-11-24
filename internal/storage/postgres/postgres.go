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
	pgx *pgx.Conn
}

type UserHistoryURL struct {
	ID       int
	CookieID string
	Hash     string
}

func New(dataSourceName string) *Database {
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

func (d *Database) GetURL(hash string) (resultURL storage.ResultURL, err error) {
	q := "SELECT * FROM urls WHERE hash = $1"
	err = d.pgx.QueryRow(context.Background(), q, hash).Scan(&resultURL.ID, &resultURL.Hash, &resultURL.Original, &resultURL.Deleted)
	if err != nil {
		fmt.Println(resultURL, err)
		return resultURL, err
	}

	//if resultURL.Original == "" {
	//	err = errors.New("not_found")
	//}

	return resultURL, err
}

func (d *Database) GetByOriginalURL(originalURL string) (hash string, err error) {
	q := "SELECT hash FROM urls WHERE original = $1"
	err = d.pgx.QueryRow(context.Background(), q, originalURL).Scan(&hash)

	if err != nil {
		return "", err
	}
	return hash, nil
}

func (d *Database) AddURL(hash, original string) error {
	q := "INSERT INTO urls (hash, original, deleted) VALUES ($1, $2, $3)"
	_, err := d.pgx.Exec(context.Background(), q, hash, original, 0)
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

	q := "INSERT INTO urls (hash, original,deleted) VALUES ($1, $2, $3)"
	_, err = tx.Prepare(context.Background(), "batch_insert", q)
	if err != nil {
		return result, err
	}

	var hash string
	for _, v := range urls {
		hash = helper.GenerateToken(8)
		_, err = tx.Exec(context.Background(), "batch_insert", hash, v.OriginalURL, 0)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.UniqueViolation {
					err = nil

				}
			}
		}
		if err == nil {
			result = append(result, storage.ResultBatchUrls{CorrelationID: v.CorrelationID, Hash: hash})
		}
	}

	if err != nil {
		_ = tx.Rollback(context.Background())
		return result, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return result, err
	}
	return result, nil
}

func (d *Database) DeleteUserBatchURL(cookieID string, hashes []string) (err error) {

	ctx := context.Background()

	tx, err := d.pgx.Begin(ctx)
	if err != nil {
		return err
	}

	b := &pgx.Batch{}

	sqlStmt := "UPDATE urls u SET deleted=1 FROM user_history_urls uh WHERE uh.hash = u.hash AND uh.hash = $1 AND uh.cookie_id = $2;"

	for _, hash := range hashes {
		b.Queue(sqlStmt, hash, cookieID)
	}

	batchResults := tx.SendBatch(ctx, b)

	var batchErr error
	for batchErr == nil {
		_, batchErr = batchResults.Exec()
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) DeleteBatchURL(hashes []string) (err error) {

	ctx := context.Background()

	tx, err := d.pgx.Begin(ctx)
	if err != nil {
		return err
	}

	b := &pgx.Batch{}

	sqlStmt := "UPDATE urls u SET deleted=1 WHERE u.hash = $1;"

	for _, hash := range hashes {
		b.Queue(sqlStmt, hash)
	}

	batchResults := tx.SendBatch(ctx, b)

	var batchErr error
	for batchErr == nil {
		_, batchErr = batchResults.Exec()
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
