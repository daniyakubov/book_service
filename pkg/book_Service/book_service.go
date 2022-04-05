package book_Service

//todo: improve error handling, use https://github.com/fiverr/go_errors/blob/master, and consult avigail

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	Action "github.com/daniyakubov/book_service/pkg/action"
	"github.com/daniyakubov/book_service/pkg/book_Service/models"
	"github.com/daniyakubov/book_service/pkg/cache"
	"github.com/daniyakubov/book_service/pkg/elastic_service"
	errors "github.com/fiverr/go_errors"
)

type BookService struct {
	booksCache     cache.Cache
	elasticHandler elastic_service.ElasticHandler
}

func NewBookService(booksCache cache.Cache, elasticHandler elastic_service.ElasticHandler) BookService {
	return BookService{
		booksCache:     booksCache,
		elasticHandler: elasticHandler,
	}
}

func (b *BookService) PutBook(req *models.Request) (string, error) {

	postBody, err := json.Marshal(req.Data)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}
	resp, err := b.elasticHandler.Put(postBody)

	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}
	var idResp models.PutBookResponse
	if err = json.Unmarshal(body, &idResp); err != nil {
		return "", errors.Wrap(err, err.Error())
	}
	err = b.booksCache.Push(req.Data.Username, "Method:Put,"+"Route:"+req.Route)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}
	return fmt.Sprintf("{id: %+v}", idResp.Id), nil

}

func (b *BookService) PostBook(req *models.Request) error {

	_, err := b.elasticHandler.Post(req.Data.Title, req.Data.Id)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}
	err = b.booksCache.Push(req.Data.Username, "Method:Post,"+"Route:"+req.Route)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}
	return nil
}

func (b *BookService) GetBook(req *models.Request) (string, error) {

	resp, err := b.elasticHandler.Get(req.Data.Id)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}
	var getResp models.GetBookResponse
	if err = json.Unmarshal(body, &getResp); err != nil {
		return "", err
	}
	src, err := json.Marshal(getResp.Source)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}

	err = b.booksCache.Push(req.Data.Username, "Method:Get,"+"Route:"+req.Route)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}
	return fmt.Sprintf("%+v", string(src)), nil

}

func (b *BookService) DeleteBook(req *models.Request) error {

	resp, err := b.elasticHandler.Delete(req.Data.Id)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}
	err = b.booksCache.Push(req.Data.Username, "Method:Delete,"+"Route:"+req.Route)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}
	return nil

}

func (b *BookService) Search(req *models.Request) (string, error) {

	resp, err := b.elasticHandler.Search(req.Data.Title, req.Data.Author, req.Data.PriceStart, req.Data.PriceEnd)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}

	var s models.SearchBookResponse
	if err := json.Unmarshal(body, &s); err != nil {
		return "", errors.Wrap(err, err.Error())
	}
	length := len(s.Hits.Hits)
	res := make([]models.Source, int(length))
	for i := 0; i < length; i++ {
		res[i] = s.Hits.Hits[i].Source
		res[i].Id = s.Hits.Hits[i].Id
	}
	postBody, err := json.Marshal(res)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}

	err = b.booksCache.Push(req.Data.Username, "Method:Get,"+"Route:"+req.Route)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}
	return fmt.Sprintf("%+v", string(postBody)), nil

}

func (b *BookService) Store(req *models.Request) (string, error) {

	resp1, resp2, err := b.elasticHandler.Store()

	defer resp1.Body.Close()

	body, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}

	var count models.StoreCount
	if err := json.Unmarshal(body, &count); err != nil {
		return "", errors.Wrap(err, err.Error())
	}

	body2, err := ioutil.ReadAll(resp2.Body)

	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}

	defer resp2.Body.Close()

	var distinctAut models.StoreDistinctAuthors
	if err := json.Unmarshal(body2, &distinctAut); err != nil {
		return "", errors.Wrap(err, err.Error())
	}
	err = b.booksCache.Push(req.Data.Username, "method:Get,"+"route:"+req.Route)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}
	return fmt.Sprintf("{books_num: %d, distinct_authors_num: %d}", count.Count, distinctAut.Hits.Total.Value), nil
}

func (b *BookService) Activity(username string) (string, error) {

	actions, err := b.booksCache.Get(username)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}
	var res []Action.Action = make([]Action.Action, int(len(actions)))

	for i := 0; i < len(actions); i++ {
		s := strings.Split(actions[i], ",")
		method := strings.Split(s[0], ":")[1]
		route := strings.Split(s[1], ":")[1]
		res[i].Method = method
		res[i].Route = route
	}

	postBody, err := json.Marshal(res)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}

	return fmt.Sprintf("%+v", string(postBody)), nil

}
