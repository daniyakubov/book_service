package main

//todo: the response shouldn't be what we get from elastic, but only the necessary fields.  it should be a json object

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/daniyakubov/book_service/pkg/ElasticService"
	"github.com/daniyakubov/book_service/pkg/cache"
)

type PutBookHit struct {
	Title  string `json:"title"`
	Author string `json:"authors_name"` //todo: make author's name be authors_name

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

type PutBookResponse struct {
	Id string `json:"_id"`
}

type GetBookResponse struct {
	Source struct {
		Title       string  `json:"title"`
		Price       float64 `json:"price"`
		AuthorsName string  `json:"authorsName"`
		Available   bool    `json:"available"`
		Date        string  `json:"date"`
	} `json:"_source"`
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

type BookService struct {
	client         http.Client
	booksCache     cache.BooksCache
	elasticHandler ElasticService.ElasticHandler
}

func NewBookService(client http.Client, booksCache cache.BooksCache, elasticHandler ElasticService.ElasticHandler) BookService {
	return BookService{
		client:         client,
		booksCache:     booksCache,
		elasticHandler: elasticHandler,
	}
}
func (b *BookService) putBook(w http.ResponseWriter, req *http.Request) {
	var hit PutBookHit
	err := json.NewDecoder(req.Body).Decode(&hit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	postBody, _ := json.Marshal(hit)
	resp, err := b.elasticHandler.Put(postBody)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	var idResp PutBookResponse
	if err := json.Unmarshal(body, &idResp); err != nil {
		fmt.Println("Can not unmarshal JSON")
		return
	}
	//var actions Actions

	fmt.Fprintf(w, "id: %+v", idResp.Id)

}

func (b *BookService) postBook(w http.ResponseWriter, req *http.Request) {
	var hit PostBookHit
	err := json.NewDecoder(req.Body).Decode(&hit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := b.elasticHandler.Post(hit.Title, hit.Id)
	if err != nil {
		fmt.Errorf("Error %s", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "result: %+v", string(body))
}

func (b *BookService) getBook(w http.ResponseWriter, req *http.Request) {
	var hit GetBookHit
	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hit.Id = req.FormValue("id")
	hit.Username = req.FormValue("username")
	resp, err := b.elasticHandler.Get(hit.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var getResp GetBookResponse
	if err := json.Unmarshal(body, &getResp); err != nil {
		fmt.Println("Can not unmarshal JSON")
		return
	}

	fmt.Fprintf(w, " %+v", getResp)
}

func (b *BookService) deleteBook(w http.ResponseWriter, req *http.Request) {
	var hit DeleteBookHit
	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hit.Id = req.FormValue("id")
	hit.Username = req.FormValue("username")

	resp, errD := b.elasticHandler.Delete(hit.Id)
	if errD != nil {
		fmt.Println(errD)
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, "response: %+v", string(respBody))

}

// handler -> client to elastic db -> parse results -> process them -> create response
func (b *BookService) book(w http.ResponseWriter, req *http.Request) {

	if req.Method == "PUT" {
		b.putBook(w, req)
	} else if req.Method == "POST" {
		b.postBook(w, req)

	} else if req.Method == "GET" {
		b.getBook(w, req)
	} else if req.Method == "DELETE" {
		b.deleteBook(w, req)
	}

}

func (b *BookService) search(w http.ResponseWriter, req *http.Request) { //todo: make the price range to be one field, and not all fields must be included
	if req.Method == "GET" {
		var hit GetSearchHit
		err := req.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		hit.Title = req.FormValue("title")
		hit.Author = req.FormValue("author_name")
		hit.PriceStart, _ = strconv.ParseFloat(req.FormValue("price_start"), 8)
		hit.PriceEnd, _ = strconv.ParseFloat(req.FormValue("price_end"), 32)
		hit.Username = req.FormValue("username")

		resp, err := b.elasticHandler.Search(hit.Title, hit.Author, hit.PriceStart, hit.PriceEnd)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, "book: %+v", string(body))
	}

}

func (b *BookService) store(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		err := req.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp1, resp2, err := b.elasticHandler.Store()

		defer resp1.Body.Close()

		body, err := ioutil.ReadAll(resp1.Body)
		if err != nil {
			return
		}

		var count StoreCount
		if err := json.Unmarshal(body, &count); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		body2, err := ioutil.ReadAll(resp2.Body)

		if err != nil {
			fmt.Println(err)
			return
		}

		defer resp2.Body.Close()

		var distinctAut StoreDistinctAuthors
		if err := json.Unmarshal(body2, &distinctAut); err != nil {
			fmt.Println("Can not unmarshal JSON")
			return
		}
		fmt.Fprintf(w, "number of books: %d, number of distinct authors: %d", count.Count, distinctAut.Hits.Total.Value)
	}
}

func (b *BookService) activity(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, req.FormValue("title"))
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {
	client := http.Client{Timeout: time.Duration(1) * time.Second}
	eHandler := ElasticService.NewElasticHandler("http://es-search-7.fiverrdev.com:9200/books/", client)
	var b BookService = NewBookService(client, cache.NewRedisCache("localhost:6379", 0, 0), eHandler)
	http.HandleFunc("/book", b.book)
	http.HandleFunc("/search", b.search)
	http.HandleFunc("/store", b.store)
	http.HandleFunc("/activity", b.activity)

	http.ListenAndServe(":8080", nil)
}
