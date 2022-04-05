package main

import (
	"net/http"
	"time"

	"github.com/daniyakubov/book_service/pkg/http_service"

	"gopkg.in/redis.v5"

	"github.com/daniyakubov/book_service/pkg/book_Service"
	"github.com/daniyakubov/book_service/pkg/cache"
	"github.com/daniyakubov/book_service/pkg/consts"
	"github.com/daniyakubov/book_service/pkg/elastic_service"
)

func main() {
	client := http.Client{Timeout: time.Duration(consts.ClientTimeOut) * time.Second}
	eHandler := elastic_service.NewElasticHandler(consts.BooksUrl, &client, consts.MaxQueryResults)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     consts.Host,
		Password: consts.Password,
		DB:       consts.Db,
	})
	bookService := book_Service.NewBookService(cache.NewRedisCache(consts.Host, consts.Db, consts.Expiration, consts.MaxActions, redisClient), eHandler)
	httpHandler := http_service.NewHttpHandler(client, bookService)

	http.HandleFunc("/book", httpHandler.Book)
	http.HandleFunc("/search", httpHandler.Search)
	http.HandleFunc("/store", httpHandler.Store)
	http.HandleFunc("/activity", httpHandler.Activity)

	err := http.ListenAndServe(consts.HttpAddress, nil)
	if err != nil {
		panic(err)
	}
}
