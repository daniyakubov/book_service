package http_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/daniyakubov/book_service/pkg/book_Service"
	"github.com/daniyakubov/book_service/pkg/book_Service/models"
)

type HttpHandler struct {
	client      http.Client
	bookService book_Service.BookService
}

func NewHttpHandler(client http.Client, bookService book_Service.BookService) HttpHandler {
	return HttpHandler{
		client:      client,
		bookService: bookService,
	}
}

func (h *HttpHandler) PutBook(w http.ResponseWriter, req *http.Request) {
	var hit models.Hit
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &hit); err != nil {
		fmt.Println("Can not unmarshal JSON")
		return
	}
	r := models.NewRequest(&hit, req.URL.Path)
	s, err := h.bookService.PutBook(&r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, s)

}

func (h *HttpHandler) PostBook(w http.ResponseWriter, req *http.Request) {
	var hit models.Hit
	err := json.NewDecoder(req.Body).Decode(&hit) //todo: change to unmarshal
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	r := models.NewRequest(&hit, req.URL.Path)
	err = h.bookService.PostBook(&r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *HttpHandler) GetBook(w http.ResponseWriter, req *http.Request) {
	var hit models.Hit
	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hit.Id = req.FormValue("id")
	hit.Username = req.FormValue("username")

	r := models.NewRequest(&hit, req.URL.Path)

	s, err := h.bookService.GetBook(&r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, s)

}

func (h *HttpHandler) DeleteBook(w http.ResponseWriter, req *http.Request) {
	var hit models.Hit
	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hit.Id = req.FormValue("id")
	hit.Username = req.FormValue("username")

	r := models.NewRequest(&hit, req.URL.Path)

	err = h.bookService.DeleteBook(&r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *HttpHandler) Book(w http.ResponseWriter, req *http.Request) {

	if req.Method == "PUT" {
		h.PutBook(w, req)
	} else if req.Method == "POST" {
		h.PostBook(w, req)

	} else if req.Method == "GET" {
		h.GetBook(w, req)
	} else if req.Method == "DELETE" {
		h.DeleteBook(w, req)
	}

}

func (h *HttpHandler) Search(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		var hit models.Hit
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
		r := models.NewRequest(&hit, req.URL.Path)

		s, err := h.bookService.Search(&r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, s)
	}

}

func (h *HttpHandler) Store(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		err := req.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var hit models.Hit
		hit.Username = req.FormValue("username")

		r := models.NewRequest(&hit, req.URL.Path)
		s, err := h.bookService.Store(&r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, s)
	}

}

func (h *HttpHandler) Activity(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := req.FormValue("username")
	s := h.bookService.Activity(user)
	fmt.Fprintf(w, s)

}
