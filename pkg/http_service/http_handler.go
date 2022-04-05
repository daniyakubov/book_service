package http_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/daniyakubov/book_service/pkg/consts"

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
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &hit); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
	hit.Id = req.FormValue(consts.Id)

	hit.Username = req.FormValue(consts.UserName)

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

	hit.Id = req.FormValue(consts.Id)
	hit.Username = req.FormValue(consts.UserName)

	r := models.NewRequest(&hit, req.URL.Path)

	err = h.bookService.DeleteBook(&r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *HttpHandler) Book(w http.ResponseWriter, req *http.Request) {

	if req.Method == consts.PutMethod {
		h.PutBook(w, req)
	} else if req.Method == consts.PostMethod {
		h.PostBook(w, req)

	} else if req.Method == consts.GetMethod {
		h.GetBook(w, req)
	} else if req.Method == consts.DeleteMethod {
		h.DeleteBook(w, req)
	}

}

func (h *HttpHandler) Search(w http.ResponseWriter, req *http.Request) {
	if req.Method == consts.GetMethod {
		var hit models.Hit
		err := req.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		hit.Title = req.FormValue(consts.Title)
		hit.Author = req.FormValue(consts.Author)
		sRange := req.FormValue(consts.PriceRange)
		if sRange == "" {
			hit.PriceStart = 0
			hit.PriceEnd = 0
		} else {
			priceRange := strings.Split(sRange, "-")

			hit.PriceStart, _ = strconv.ParseFloat(priceRange[0], 32)
			hit.PriceEnd, _ = strconv.ParseFloat(priceRange[1], 32)
		}
		hit.Username = req.FormValue(consts.UserName)
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
	if req.Method == consts.GetMethod {
		err := req.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var hit models.Hit
		hit.Username = req.FormValue(consts.UserName)

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
	user := req.FormValue(consts.UserName)
	s, err := h.bookService.Activity(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, s)

}
