package main

import (
	"net/http"
	"time"

	"github.com/daniyakubov/book_service/pkg/http_service"

	"gopkg.in/redis.v5"

	"github.com/daniyakubov/book_service/pkg/book_Service"
	"github.com/daniyakubov/book_service/pkg/cache"
	"github.com/daniyakubov/book_service/pkg/elastic_service"
)

func main() {
	client := http.Client{Timeout: time.Duration(10) * time.Second}
	eHandler := elastic_service.NewElasticHandler("http://es-search-7.fiverrdev.com:9200/books/", &client, 10000)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	bookService := book_Service.NewBookService(cache.NewRedisCache("localhost:6379", 0, 0, 3, redisClient), eHandler)
	httpHandler := http_service.NewHttpHandler(client, bookService)

	http.HandleFunc("/book", httpHandler.Book)
	http.HandleFunc("/search", httpHandler.Search)
	http.HandleFunc("/store", httpHandler.Store)
	http.HandleFunc("/activity", httpHandler.Activity)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
