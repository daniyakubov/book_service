package models

type PutBookResponse struct {
	Id string `json:"_id"`
}

type GetBookResponse struct {
	Source struct {
		Title     string  `json:"title"`
		Price     float64 `json:"price"`
		Author    string  `json:"author"`
		Available bool    `json:"available"`
		Date      string  `json:"date"`
	} `json:"_source"`
}

type Source struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Author      string  `json:"author"`
	Available   bool    `json:"available"`
	Date        string  `json:"date,omitempty"`
	PublishDate string  `json:"publish,date,omitempty"`
}

type RHit struct {
	Index  string  `json:"_index"`
	Type   string  `json:"_type"`
	Id     string  `json:"_id"`
	Score  float64 `json:"_score"`
	Source Source  `json:"_source"`
}

type SearchBookResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []RHit  `json:"hits"`
	} `json:"hits"`
}

type StoreDistinctAuthors struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
	} `json:"hits"`
}

type StoreCount struct {
	Count int `json:"count"`
}
