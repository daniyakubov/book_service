package models

type PutBookHit struct {
	Title  string `json:"title"`
	Author string `json:"author"`

	Price     float32 `json:"price"`
	Available bool    `json:"available"`
	Date      string  `json:"date"`
	Username  string  `json:"username"`
}

type PostBookHit struct {
	Id       string
	Title    string
	Username string
}

type GetBookHit struct {
	Id       string
	Username string
}

type DeleteBookHit struct {
	Id       string
	Username string
}

type GetSearchHit struct {
	Title      string
	Author     string
	PriceStart float64
	PriceEnd   float64
	Username   string
}
