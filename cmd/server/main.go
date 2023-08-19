package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/igoramorim/dollar-exchange-rate/internal/exchrate"
	"github.com/igoramorim/dollar-exchange-rate/internal/repository"
	"github.com/igoramorim/dollar-exchange-rate/internal/repository/sqlite"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Println(err)
	}
}

func run() error {
	db, err := sql.Open("sqlite3", "./dolar_exchange_rate.db")
	if err != nil {
		return errors.WithMessage(err, "open sqlite3 connection")
	}
	defer db.Close()

	err = setupTable(db)
	if err != nil {
		return errors.WithMessage(err, "setup sqlite3 table")
	}

	repo, err := sqlite.New(db)
	if err != nil {
		return errors.WithMessage(err, "create repository")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", handleDollarExchangeRate(repo))
	return http.ListenAndServe(":8080", mux)
}

func handleDollarExchangeRate(repo repository.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL)

		ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
		defer cancel()

		exchRate, err := exchrate.Get(ctx)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		ctx, cancel = context.WithTimeout(r.Context(), 10*time.Millisecond)
		defer cancel()

		err = repo.Save(ctx, exchRate.ExchangeRate)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(exchRate); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		return
	}
}

func setupTable(db *sql.DB) error {
	const stmt = `
		create table if not exists dolar_exchange_rate(
			rate decimal(4) not null,
		    created_at timestamp(6) default current_timestamp
		);
	`
	_, err := db.Exec(stmt)
	return err
}
