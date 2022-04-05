package models

type Hit struct {
	Id         string
	Title      string  `json:"title"`
	Author     string  `json:"author"`
	Price      float32 `json:"price"`
	Available  bool    `json:"available"`
	Date       string  `json:"date"`
	Username   string  `json:"username"`
	PriceStart float64
	PriceEnd   float64
}
