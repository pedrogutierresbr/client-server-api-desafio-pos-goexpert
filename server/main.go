package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	apiTimeout = 1000 * time.Millisecond
	dbTimeout  = 1000 * time.Millisecond
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
	dsn := "root:root@tcp(localhost:3306)/goexpert?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Dollar{})

	gormCtx, gormCancel := context.WithTimeout(context.Background(), dbTimeout)
	defer gormCancel()

	select {
	case <-gormCtx.Done():
		return nil, fmt.Errorf("timeout reached!\n%v", err)
	default:
		bidDollar := &Dollar{
			ID:  uuid.New(),
			Bid: quotation.USDBRL.Bid,
		}

		db.WithContext(gormCtx).Create(bidDollar)
		return bidDollar, nil
	}

}
