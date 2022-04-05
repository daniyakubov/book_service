package main

import (
	"net/http"
	"time"

	"gopkg.in/redis.v5"

	"github.com/daniyakubov/book_service/pkg/book_Service"
	"github.com/daniyakubov/book_service/pkg/cache"
	"github.com/daniyakubov/book_service/pkg/elastic_service"
)

func main() {
	client := http.Client{Timeout: time.Duration(10) * time.Second}
	eHandler := elastic_service.NewElasticHandler("http://es-search-7.fiverrdev.com:9200/books/", &client)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	b := book_Service.NewBookService(client, cache.NewRedisCache("localhost:6379", 0, 0, 3, redisClient), eHandler)
	http.HandleFunc("/book", b.Book)
	http.HandleFunc("/search", b.Search)
	http.HandleFunc("/store", b.Store)
	http.HandleFunc("/activity", b.Activity)

	http.ListenAndServe(":8080", nil)
}
