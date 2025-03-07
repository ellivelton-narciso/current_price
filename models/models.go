package models

import "time"

type PriceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
	Time   int64  `json:"time"`
}

type ConfigStruct struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	DBname    string `json:"dbname"`
	User      string `json:"user"`
	Pass      string `json:"pass"`
	BaseURL   string `json:"baseURL"`
	Leituras  int    `json:"leituras"`
	TimeSleep int64  `json:"timeSleep"`
}

type Historico struct {
	Id        int64           `json:"id"`
	Value     []PriceResponse `json:"value"`
	CreatedAt time.Time       `json:"created_at"`
}

type Bots struct {
	Symbol string `json:"symbol"`
	User   string `json:"user"`
}
