package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/igoramorim/dollar-exchange-rate/internal/exchrate"
	"github.com/igoramorim/dollar-exchange-rate/internal/repository/txtfile"
	"github.com/pkg/errors"
	"io"
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
	repo, err := txtfile.New("cotacao.txt")
	if err != nil {
		return err
	}
	defer repo.Close()

	ctx := context.Background()

	exchangeRate, err := getExchangeRate(ctx)
	if err != nil {
		return errors.WithMessage(err, "get exchange rate")
	}

	if err = repo.Save(ctx, exchangeRate); err != nil {
		return errors.WithMessage(err, "save exchange rate")
	}

	return nil
}

func getExchangeRate(ctx context.Context) (float64, error) {
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	const url = "http://localhost:8080/cotacao"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, errors.WithMessage(err, "create request")
	}

	log.Printf("get exchange rate on %s\n", url)

	httpRes, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, errors.WithMessage(err, "execute request")
	}

	log.Printf("response status: %d\n", httpRes.StatusCode)

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return 0, errors.WithMessage(err, "read response body")
	}

	if httpRes.StatusCode > 200 {
		return 0, errors.New(fmt.Sprintf("api response code: %d %s", httpRes.StatusCode, string(body)))
	}

	var apiRes exchrate.Response
	if err = json.Unmarshal(body, &apiRes); err != nil {
		return 0, errors.WithMessage(err, "parse api response")
	}

	log.Printf("response body: %+v\n", apiRes)

	return apiRes.ExchangeRate, nil
}
