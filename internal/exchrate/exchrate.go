package exchrate

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"strconv"
)

const url = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

func Get(ctx context.Context) (Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Response{}, err
	}

	log.Printf("exchrate: get exchange rate on %s", url)

	httpRes, err := http.DefaultClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer httpRes.Body.Close()

	log.Printf("exchrate: response status: %d\n", httpRes.StatusCode)

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return Response{}, err
	}

	if httpRes.StatusCode > 200 {
		return Response{}, errors.New(fmt.Sprintf("api response code: %d %s", httpRes.StatusCode, string(body)))
	}

	var apiRes apiResponse
	if err = json.Unmarshal(body, &apiRes); err != nil {
		return Response{}, err
	}

	exchangeRate, err := strconv.ParseFloat(apiRes.Usdbrl.Bid, 64)
	if err = json.Unmarshal(body, &apiRes); err != nil {
		return Response{}, err
	}

	log.Printf("exchrate: response: %+v\n", apiRes)

	return Response{ExchangeRate: exchangeRate}, nil
}

type Response struct {
	ExchangeRate float64 `json:"exchangeRate"`
}

type apiResponse struct {
	Usdbrl usdbrl `json:"USDBRL"`
}

type usdbrl struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}
