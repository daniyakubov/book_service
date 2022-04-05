package book_Service

//todo: improve error handling, use https://github.com/fiverr/go_errors/blob/master, and consult avigail

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/daniyakubov/book_service/pkg/book_Service/models"
	"github.com/daniyakubov/book_service/pkg/cache"
	"github.com/daniyakubov/book_service/pkg/elastic_service"
)

type BookService struct {
	client         http.Client
	booksCache     cache.Cache
	elasticHandler elastic_service.ElasticHandler
}

func NewBookService(client http.Client, booksCache cache.Cache, elasticHandler elastic_service.ElasticHandler) BookService {
	return BookService{
		client:         client,
		booksCache:     booksCache,
		elasticHandler: elasticHandler,
	}
}
func (b *BookService) PutBook(w http.ResponseWriter, req *http.Request) { //todo: seperate the http
	var hit models.PutBookHit
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &hit); err != nil {
		fmt.Println("Can not unmarshal JSON")
		return
	}

	postBody, err := json.Marshal(hit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := b.elasticHandler.Put(postBody)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	var idResp models.PutBookResponse
	if err := json.Unmarshal(body, &idResp); err != nil {
		fmt.Println("Can not unmarshal JSON")
		return
	}
	//var actions Actions

	fmt.Fprintf(w, "{id: %+v}", idResp.Id)

	b.booksCache.Push(hit.Username, "Method:Put,"+"Route:"+req.URL.Path)

}

func (b *BookService) PostBook(w http.ResponseWriter, req *http.Request) {
	var hit models.PostBookHit
	err := json.NewDecoder(req.Body).Decode(&hit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = b.elasticHandler.Post(hit.Title, hit.Id)
	if err != nil {
		fmt.Errorf("Error %s", err)
		return
	}
	b.booksCache.Push(hit.Username, "Method:Post,"+"Route:"+req.URL.Path)

}

func (b *BookService) GetBook(w http.ResponseWriter, req *http.Request) {
	var hit models.GetBookHit
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
	var getResp models.GetBookResponse
	if err := json.Unmarshal(body, &getResp); err != nil {
		fmt.Println("Can not unmarshal JSON")
		return
	}
	src, err := json.Marshal(getResp.Source)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, " %+v", string(src))
	b.booksCache.Push(hit.Username, "Method:Get,"+"Route:"+req.URL.Path)

}

func (b *BookService) DeleteBook(w http.ResponseWriter, req *http.Request) {
	var hit models.DeleteBookHit
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

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	b.booksCache.Push(hit.Username, "Method:Delete,"+"Route:"+req.URL.Path)

}

// handler -> client to elastic db -> parse results -> process them -> create response
func (b *BookService) Book(w http.ResponseWriter, req *http.Request) {

	if req.Method == "PUT" {
		b.PutBook(w, req)
	} else if req.Method == "POST" {
		b.PostBook(w, req)

	} else if req.Method == "GET" {
		b.GetBook(w, req)
	} else if req.Method == "DELETE" {
		b.DeleteBook(w, req)
	}

}

func (b *BookService) Search(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		var hit models.GetSearchHit
		err := req.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		hit.Title = req.FormValue("title")
		hit.Author = req.FormValue("author")
		sRange := req.FormValue("price_range")
		if sRange == "" {
			hit.PriceStart = 0
			hit.PriceEnd = 0
		} else {
			priceRange := strings.Split(sRange, "-")

			hit.PriceStart, _ = strconv.ParseFloat(priceRange[0], 32)
			hit.PriceEnd, _ = strconv.ParseFloat(priceRange[1], 32)
		}
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

		var s models.SearchBookResponse
		if err := json.Unmarshal(body, &s); err != nil {
			fmt.Println("Can not unmarshal JSON")
			return
		}
		length := len(s.Hits.Hits)
		var res []models.Source = make([]models.Source, int(length))
		for i := 0; i < length; i++ {
			res[i] = s.Hits.Hits[i].Source
		}
		postBody, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, " %+v", string(postBody))
		b.booksCache.Push(hit.Username, "Method:Get,"+"Route:"+req.URL.Path)

	}

}

func (b *BookService) Store(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		err := req.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//user := req.FormValue("username")
		//actions := b.booksCache.Get(user)
		//newActions := BookCache.
		err = req.ParseForm()
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

		var count models.StoreCount
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

		var distinctAut models.StoreDistinctAuthors
		if err := json.Unmarshal(body2, &distinctAut); err != nil {
			fmt.Println("Can not unmarshal JSON")
			return
		}
		fmt.Fprintf(w, "{books_num: %d, distinct_authors_num: %d}", count.Count, distinctAut.Hits.Total.Value)
		user := req.FormValue("username")
		b.booksCache.Push(user, "Method:Get,"+"Route:"+req.URL.Path)

	}
}

func (b *BookService) Activity(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := req.FormValue("username")

	actions := b.booksCache.Get(user)
	fmt.Fprintf(w, actions)

}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v", name, h)
		}
	}

}
