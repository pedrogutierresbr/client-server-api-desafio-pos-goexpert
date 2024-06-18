package main

import "time"

type quotationDetails struct {
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

const (
	apiURL     = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	apiTimeout = 200 * time.Millisecond
	dbTimeout  = 10 * time.Microsecond
)

func main() {

}
