package main

import (
	"net/http"
	"time"

	"github.com/daniyakubov/book_service/pkg/BookService"
	"github.com/daniyakubov/book_service/pkg/ElasticService"
	"github.com/daniyakubov/book_service/pkg/cache"
)

func main() {
	client := http.Client{Timeout: time.Duration(10) * time.Second}
	eHandler := ElasticService.NewElasticHandler("http://es-search-7.fiverrdev.com:9200/books/", client)
	var b BookService.BookService = BookService.NewBookService(client, cache.NewRedisCache("localhost:6379", 0, 0), eHandler)
	http.HandleFunc("/book", b.Book)
	http.HandleFunc("/search", b.Search)
	http.HandleFunc("/store", b.Store)
	http.HandleFunc("/activity", b.Activity)

	http.ListenAndServe(":8080", nil)
}
