package elastic_service

import (
	"bytes"
	"fmt"
	"net/http"

	errors "github.com/fiverr/go_errors"
)

type ElasticHandler struct {
	Url          string
	Client       *http.Client
	maxSizeQuery int
}

func NewElasticHandler(url string, client *http.Client, maxSizeQuery int) ElasticHandler {
	return ElasticHandler{
		url,
		client,
		maxSizeQuery,
	}
}

func (e *ElasticHandler) Post(title string, id string) (resp *http.Response, err error) {
	s := fmt.Sprintf(`{"doc": {"title": "%s"}}`, title)
	myJson := bytes.NewBuffer([]byte(s))

	resp, err = e.Client.Post(e.Url+"_update/"+id, "application/json", myJson)
	if err != nil {
		return resp, errors.Wrap(err, err.Error())
	}
	return resp, err
}

func (e *ElasticHandler) Put(postBody []byte) (resp *http.Response, err error) {
	resp, err = e.Client.Post(e.Url+"_doc/", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return resp, errors.Wrap(err, err.Error())
	}
	return resp, err
}
func (e *ElasticHandler) Get(id string) (resp *http.Response, err error) {
	resp, err = e.Client.Get(e.Url + "_doc/" + id)
	if err != nil {
		return resp, errors.Wrap(err, err.Error())
	}
	return resp, err
}

func (e *ElasticHandler) Delete(id string) (resp *http.Response, err error) {
	req, err := http.NewRequest("DELETE", e.Url+"_doc/"+id, nil)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}
	resp, err = e.Client.Do(req)
	if err != nil {
		return resp, errors.Wrap(err, err.Error())
	}
	return resp, err

}

func (e *ElasticHandler) Search(title string, author string, priceStart float64, priceEnd float64) (resp *http.Response, err error) {
	s := ""
	if title == "" && author == "" && priceEnd == 0 {
		s = fmt.Sprintf(`{"size": %d, "query": {"match_all": {}}}`, e.maxSizeQuery)
	} else if title == "" && author == "" {
		s = fmt.Sprintf(`{"size": %d, "query": {"range": {"price": {"gte": %f, "lte": %f}}}}`, e.maxSizeQuery, priceStart, priceEnd)
	} else if title == "" && priceEnd == 0 {
		s = fmt.Sprintf(`{"size": %d, "query": {"constant_score": {"filter": {"bool": {"must":[{"match": {"author": "%s"}}]}}}}}`, e.maxSizeQuery, author)
	} else if author == "" && priceEnd == 0 {
		s = fmt.Sprintf(`{"size": %d, "query": {"constant_score": {"filter": {"bool": {"must":[{"match": {"title": "%s"}}]}}}}}`, e.maxSizeQuery, title)
	} else if priceEnd == 0 {
		s = fmt.Sprintf(`{"size": %d, "query": {"constant_score": {"filter": {"bool": {"must":[{"match": {"title": "%s"}},{"match": {"author": "%s"}}]}}}}}`, e.maxSizeQuery, title, author)
	} else if title == "" {
		s = fmt.Sprintf(`{"size": %d, "query": {"constant_score": {"filter": {"bool": {"must":[{"match": {"author": "%s"}},{"range": {"price": {"gte": %f, "lte": %f} }}]}}}}}`, e.maxSizeQuery, author, priceStart, priceEnd)
	} else if author == "" {
		s = fmt.Sprintf(`{"size": %d, "query": {"constant_score": {"filter": {"bool": {"must":[{"match": {"title": "%s"}},{"range": {"price": {"gte": %f, "lte": %f} }}]}}}}}`, e.maxSizeQuery, title, priceStart, priceEnd)
	} else {
		s = fmt.Sprintf(`{"size": %d, "query": {"constant_score": {"filter": {"bool": {"must":[{"match": {"title": "%s"}},{"match": {"author": "%s"}},{"range": {"price": {"gte": %f, "lte": %f} }}]}}}}}`, e.maxSizeQuery, title, author, priceStart, priceEnd)
	}

	myJson := bytes.NewBuffer([]byte(s))

	req, err := http.NewRequest("GET", e.Url+"_search/", myJson)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = e.Client.Do(req)
	if err != nil {
		return resp, errors.Wrap(err, err.Error())
	}
	return resp, err
}

func (e *ElasticHandler) Store() (resp1 *http.Response, resp2 *http.Response, err error) {
	s1 := fmt.Sprintf(`{"query": {"match_all": {}}}`)
	myJson := bytes.NewBuffer([]byte(s1))

	req, err := http.NewRequest("GET", e.Url+"_count/", myJson)
	if err != nil {
		return nil, nil, errors.Wrap(err, err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	resp1, err = e.Client.Do(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, err.Error())
	}

	s2 := fmt.Sprintf(`{"aggs" : {"authors_count" : {"cardinality" : {"field" : "authorsName.keyword"}}}}`)
	myJson2 := bytes.NewBuffer([]byte(s2))

	resp2, err = e.Client.Post(e.Url+"_search/", "application/json", myJson2)
	if err != nil {
		return nil, nil, errors.Wrap(err, err.Error())
	}
	return resp1, resp2, nil
}
