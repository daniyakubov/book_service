package book_Service

//todo: improve error handling, use https://github.com/fiverr/go_errors/blob/master, and consult avigail

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/daniyakubov/book_service/pkg/book_Service/models"
	"github.com/daniyakubov/book_service/pkg/cache"
	"github.com/daniyakubov/book_service/pkg/elastic_service"
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
		return "", err
	}
	resp, err := b.elasticHandler.Put(postBody)

	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	var idResp models.PutBookResponse
	if err := json.Unmarshal(body, &idResp); err != nil {
		return "", err
	}
	b.booksCache.Push(req.Data.Username, "Method:Put,"+"Route:"+req.Route)
	return fmt.Sprintf("{id: %+v}", idResp.Id), nil

}

func (b *BookService) PostBook(req *models.Request) error {

	_, err := b.elasticHandler.Post(req.Data.Title, req.Data.Id)
	if err != nil {
		return err
	}
	b.booksCache.Push(req.Data.Username, "Method:Post,"+"Route:"+req.Route)
	return nil
}

func (b *BookService) GetBook(req *models.Request) (string, error) {

	resp, err := b.elasticHandler.Get(req.Data.Id)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var getResp models.GetBookResponse
	if err := json.Unmarshal(body, &getResp); err != nil {
		return "", err
	}
	src, err := json.Marshal(getResp.Source)
	if err != nil {
		return "", err
	}

	b.booksCache.Push(req.Data.Username, "Method:Get,"+"Route:"+req.Route)
	return fmt.Sprintf("%+v", string(src)), nil

}

func (b *BookService) DeleteBook(req *models.Request) error {

	resp, err := b.elasticHandler.Delete(req.Data.Id)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	b.booksCache.Push(req.Data.Username, "Method:Delete,"+"Route:"+req.Route)
	return nil

}

func (b *BookService) Search(req *models.Request) (string, error) {

	resp, err := b.elasticHandler.Search(req.Data.Title, req.Data.Author, req.Data.PriceStart, req.Data.PriceEnd)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var s models.SearchBookResponse
	if err := json.Unmarshal(body, &s); err != nil {
		return "", err
	}
	length := len(s.Hits.Hits)
	var res []models.Source = make([]models.Source, int(length))
	for i := 0; i < length; i++ {
		res[i] = s.Hits.Hits[i].Source
	}
	postBody, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	b.booksCache.Push(req.Data.Username, "Method:Get,"+"Route:"+req.Route)
	return fmt.Sprintf("%+v", string(postBody)), nil

}

func (b *BookService) Store(req *models.Request) (string, error) {

	resp1, resp2, err := b.elasticHandler.Store()

	defer resp1.Body.Close()

	body, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		return "", err
	}

	var count models.StoreCount
	if err := json.Unmarshal(body, &count); err != nil {
		return "", err
	}

	body2, err := ioutil.ReadAll(resp2.Body)

	if err != nil {
		return "", err
	}

	defer resp2.Body.Close()

	var distinctAut models.StoreDistinctAuthors
	if err := json.Unmarshal(body2, &distinctAut); err != nil {
		return "", err
	}
	b.booksCache.Push(req.Data.Username, "Method:Get,"+"Route:"+req.Route)
	return fmt.Sprintf("{books_num: %d, distinct_authors_num: %d}", count.Count, distinctAut.Hits.Total.Value), nil
}

func (b *BookService) Activity(username string) string {

	actions := b.booksCache.Get(username)
	return actions

}
