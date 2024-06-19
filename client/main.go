package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type QuotationValue struct {
	Bid string `json:"bid"`
}

const (
	serverURL     = "http://localhost:8080/cotacao"
	clientTimeout = 300 * time.Millisecond
	outputFileName
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	quotation := getQuotation(ctx)

	fmt.Printf("valor da quotation %v", quotation)
}

func getQuotation(ctx context.Context) QuotationValue {
	req, err := http.NewRequestWithContext(ctx, "GET", serverURL, nil)
	if err != nil {
		log.Panic(err)
	}

	select {
	case <-ctx.Done():
		log.Panic("Timeout reached.")
		return QuotationValue{}
	default:
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Panic("Ocorreu um erro ao tentar ler o response da consulta da cotação do dolar")
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		var quotationDolar QuotationValue
		err = json.Unmarshal(body, &quotationDolar)
		if err != nil {
			panic(err)
		}

		return quotationDolar
	}
}
