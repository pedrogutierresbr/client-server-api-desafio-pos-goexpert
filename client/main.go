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
	serverURL      = "http://localhost:8080/cotacao"
	clientTimeout  = 300 * time.Millisecond
	outputFileName = "cotacao.txt"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	quotation, err := GetQuotation(ctx)
	if err != nil {
		log.Fatalf("Erro ao obter cotação: %v", err)
	}

	fmt.Printf("valor da quotation %v", quotation)
}

func GetQuotation(ctx context.Context) (QuotationValue, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", serverURL, nil)
	if err != nil {
		return QuotationValue{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return QuotationValue{}, fmt.Errorf("Timeout reached!")
		default:
			return QuotationValue{}, fmt.Errorf("Erro ao fazer a requisição: %v", err)
		}
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return QuotationValue{}, fmt.Errorf("Status code inesperado: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return QuotationValue{}, err
	}

	var quotationDolar QuotationValue
	err = json.Unmarshal(body, &quotationDolar)
	if err != nil {
		panic(err)
	}

	return quotationDolar, nil
}
