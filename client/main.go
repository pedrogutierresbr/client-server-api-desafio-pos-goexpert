package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type QuotationValue struct {
	Bid string `json:"bid"`
}

const (
	clientTimeout  = 1000 * time.Millisecond
	outputFileName = "cotacao.txt"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	quotation, err := GetQuotation(ctx)
	if err != nil {
		log.Fatalf("Erro ao obter cotação: %v", err)
	}

	err = SaveQuotationFile(quotation)
	if err != nil {
		log.Fatalf("Erro ao salvar valor no arquivo: %v", err)
	}

	fmt.Printf("Obrigado por usar o meu programa!")
}

func GetQuotation(ctx context.Context) (QuotationValue, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return QuotationValue{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return QuotationValue{}, fmt.Errorf("timeout reached!\n%v", err)
		default:
			return QuotationValue{}, fmt.Errorf("erro ao fazer a requisição!\n%v", err)
		}
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return QuotationValue{}, fmt.Errorf("status code inesperado: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return QuotationValue{}, err
	}

	var quotationDolar QuotationValue
	err = json.Unmarshal(body, &quotationDolar)
	if err != nil {
		fmt.Printf("Erro ao parsear valor do câmbio!\n%v", err)
	}

	return quotationDolar, nil
}

func SaveQuotationFile(quotation QuotationValue) error {
	quotationFile, err := os.Create("cotacao.txt")
	if err != nil {
		return fmt.Errorf("erro ao criar o arquivo de cotação: %v", err)
	}
	defer quotationFile.Close()

	_, err = quotationFile.WriteString("Dólar: " + quotation.Bid + "\n")
	if err != nil {
		return fmt.Errorf("erro ao escrever no arquivo de cotação: %v", err)
	}

	return nil
}
