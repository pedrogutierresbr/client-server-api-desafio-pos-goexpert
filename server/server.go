package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	uuid "github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type QuotationDetails struct {
	ID         int    `json:"id"`
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

type QuotationResponse struct {
	USDBRL QuotationDetails
}

type Dollar struct {
	ID  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Bid string    `json:"bid"`
}

const (
	apiURL     = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	apiTimeout = 200 * time.Millisecond
	dbTimeout  = 10 * time.Millisecond
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", QuotationHandler)
	http.ListenAndServe(":8080", mux)
}

func QuotationHandler(w http.ResponseWriter, r *http.Request) {
	quotation, err := GetDollarQuotation()
	if err != nil {
		log.Printf("Erro ao obter cotação: %v", err)
		http.Error(w, "Erro ao obter cotação", http.StatusInternalServerError)
		return
	}

	bid, err := SaveQuotationInDB(*quotation)
	if err != nil {
		log.Printf("Erro ao salvar cotação no banco de dados: %v", err)
		http.Error(w, "Erro ao salvar cotação no banco de dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(bid)
	if err != nil {
		log.Printf("Erro ao codificar resposta JSON: %v", err)
		http.Error(w, "Erro ao codificar resposta JSON", http.StatusInternalServerError)
	}
}

func GetDollarQuotation() (*QuotationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request de consulta de cotação do dólar %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao realizar consulta a API! %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("erro na leitura da response da API! %v", err)
	}

	var quotationResponse QuotationResponse
	err = json.Unmarshal(body, &quotationResponse)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter response da API! %v", err)
	}

	return &quotationResponse, nil
}

func SaveQuotationInDB(quotation QuotationResponse) (*Dollar, error) {
	dsn := "root:root@tcp(localhost:3306)/goexpert?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %v", err)
	}
	db.AutoMigrate(&Dollar{})

	gormCtx, gormCancel := context.WithTimeout(context.Background(), dbTimeout)
	defer gormCancel()

	select {
	case <-gormCtx.Done():
		return nil, fmt.Errorf("erro ao salvar cotação no banco de dados: timeout reachead!\n%v", gormCtx.Err())
	default:
		bidDollar := &Dollar{
			ID:  uuid.New(),
			Bid: quotation.USDBRL.Bid,
		}

		err := db.WithContext(gormCtx).Create(bidDollar)
		if err != nil {
			return nil, fmt.Errorf("erro ao salvar cotação no banco de dados: %v", err)
		}
		return bidDollar, nil
	}

}
