package models

type Stock struct {
	StockID int64  `json:"stockid"`
	Name    string `json:"name"`
	Price   int    `json:"price"`
	Company string `json:"company"`
}
