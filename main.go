package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type PutBookHit struct {
	Title     string
	Author    string
	Price     float32
	Available bool
	Date      string
	Username  string
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

func putBook(w http.ResponseWriter, req *http.Request, client http.Client) {
	var hit PutBookHit
	err := json.NewDecoder(req.Body).Decode(&hit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	postBody, _ := json.Marshal(hit)
	resp, err := client.Post("http://es-search-7.fiverrdev.com:9200/books/_doc/", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "book: %+v", string(body))

}

func postBook(w http.ResponseWriter, req *http.Request, client http.Client) {
	var hit PostBookHit
	err := json.NewDecoder(req.Body).Decode(&hit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := fmt.Sprintf(`{"doc": {"title": "%s"}}`, hit.Title)
	myJson := bytes.NewBuffer([]byte(s))

	resp, err := client.Post("http://es-search-7.fiverrdev.com:9200/books/_update/"+hit.Id, "application/json", myJson)

	if err != nil {
		fmt.Errorf("Error %s", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "book: %+v", string(body))
}

func getBook(w http.ResponseWriter, req *http.Request, client http.Client) {
	var hit GetBookHit
	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hit.Id = req.FormValue("id")
	hit.Username = req.FormValue("username")
	resp, err := client.Get("http://es-search-7.fiverrdev.com:9200/books/_doc/" + hit.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "book: %+v", string(body))
}

func deleteBook(w http.ResponseWriter, req *http.Request, client http.Client) {
	var hit DeleteBookHit
	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hit.Id = req.FormValue("id")
	hit.Username = req.FormValue("username")

	reqD, errD := http.NewRequest("DELETE", "http://es-search-7.fiverrdev.com:9200/books/_doc/"+hit.Id, nil)
	if errD != nil {
		fmt.Println(errD)
		return
	}

	resp, errD := client.Do(reqD)
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
func book(w http.ResponseWriter, req *http.Request) {
	client := http.Client{Timeout: time.Duration(1) * time.Second}
	if req.Method == "PUT" {
		putBook(w, req, client)
	} else if req.Method == "POST" {
		postBook(w, req, client)

	} else if req.Method == "GET" {
		getBook(w, req, client)
	} else if req.Method == "DELETE" {
		deleteBook(w, req, client)
	}

}

func search(w http.ResponseWriter, req *http.Request) {
	client := http.Client{Timeout: time.Duration(1) * time.Second}
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

		//{"query": {"constant_score": {"filter": {"bool": {"must": [{"match": {"title": "Harry Potter"}},{"match": {"authorsName": "Anna Byers"}},{"range": {"price": {"gte": 50} }}]}}}}}
		s := fmt.Sprintf(`{"query": {"constant_score": {"filter": {"bool": {"must":[{"match": {"title": "%s"}},{"match": {"authorsName": "%s"}},{"range": {"price": {"gte": %f, "lte": %f} }}]}}}}}`, hit.Title, hit.Author, hit.PriceStart, hit.PriceEnd)
		myJson := bytes.NewBuffer([]byte(s))

		req, err := http.NewRequest("GET", "http://es-search-7.fiverrdev.com:9200/books/_search/", myJson)
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, "book: %+v", string(body))
	}

}

func store(w http.ResponseWriter, req *http.Request) {
	client := http.Client{Timeout: time.Duration(1) * time.Second}
	if req.Method == "GET" {
		err := req.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s1 := fmt.Sprintf(`{"query": {"match_all": {}}}`)
		myJson := bytes.NewBuffer([]byte(s1))

		req, err := http.NewRequest("GET", "http://es-search-7.fiverrdev.com:9200/books/_count/", myJson)
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		s2 := fmt.Sprintf(`{"aggs" : {"authors_count" : {"cardinality" : {"field" : "authorsName.keyword"}}}}`)
		myJson2 := bytes.NewBuffer([]byte(s2))

		resp, err = client.Post("http://es-search-7.fiverrdev.com:9200/books/_search/", "application/json", myJson2)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		body2, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Println(err)
			return
		}

		defer resp.Body.Close()

		fmt.Fprintf(w, "book: %+v", string(body)+"dfsaf"+string(body2))
	}
}

func activity(w http.ResponseWriter, req *http.Request) {
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

	http.HandleFunc("/book", book)
	http.HandleFunc("/search", search)
	http.HandleFunc("/store", store)
	http.HandleFunc("/activity", activity)

	http.ListenAndServe(":8080", nil)
}
