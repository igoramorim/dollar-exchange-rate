package sqlite

import (
	"context"
	"database/sql"
	"github.com/igoramorim/dollar-exchange-rate/internal/repository"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func New(db *sql.DB) (*SQLite, error) {
	log.Println("initializing sqlite repository")
	return &SQLite{
		db: db,
	}, nil
}

var _ repository.Repository = (*SQLite)(nil)

type SQLite struct {
	db *sql.DB
}

func (r *SQLite) Save(ctx context.Context, exchangeRate float64) error {
	log.Println("sqlite: saving", exchangeRate)

	stmt, err := r.db.Prepare("insert into dolar_exchange_rate(rate) values(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, exchangeRate)
	if err != nil {
		return err
	}

	return nil
}
