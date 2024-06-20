package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	uuid "github.com/google/uuid"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
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
	ID  uuid.UUID `gorm:"type:uuid;primaryKey;" json:"-"`
	Bid string    `json:"bid"`
}

const (
	apiURL     = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	apiTimeout = 600 * time.Millisecond
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
		panic(err)
	}

	bid, err := SaveQuotationInDB(*quotation)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bid)
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
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Dollar{})

	dbCtx, dbCancel := context.WithTimeout(context.Background(), dbTimeout)
	defer dbCancel()

	select {
	case <-dbCtx.Done():
		fmt.Println("DB Timeout!")
		return nil, dbCtx.Err()
	default:
		bidDollar := &Dollar{
			ID:  uuid.New(),
			Bid: quotation.USDBRL.Bid,
		}

		db.WithContext(dbCtx).Create(&bidDollar)
		return bidDollar, nil
	}
}
