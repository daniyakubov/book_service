package ElasticService

import (
	"bytes"
	"fmt"
	"net/http"
)

type ElasticHandler struct {
	Url    string
	Client http.Client
}

func NewElasticHandler(url string, client http.Client) ElasticHandler {
	return ElasticHandler{
		url,
		client,
	}
}

func (e *ElasticHandler) Post(title string, id string) (resp *http.Response, err error) {
	s := fmt.Sprintf(`{"doc": {"title": "%s"}}`, title)
	myJson := bytes.NewBuffer([]byte(s))

	resp, err = e.Client.Post(e.Url+"_update/"+id, "application/json", myJson)
	return resp, err
}

func (e *ElasticHandler) Put(postBody []byte) (resp *http.Response, err error) {
	return e.Client.Post(e.Url+"_doc/", "application/json", bytes.NewBuffer(postBody))

}
func (e *ElasticHandler) Get(id string) (resp *http.Response, err error) {
	return e.Client.Get(e.Url + "_doc/" + id)

}

func (e *ElasticHandler) Delete(id string) (resp *http.Response, err error) {
	req, err := http.NewRequest("DELETE", e.Url+"_doc/"+id, nil)
	if err != nil {
		return nil, err
	}

	return e.Client.Do(req)

}

func (e *ElasticHandler) Search(title string, author string, priceStart float64, piceEnd float64) (resp *http.Response, err error) {
	s := fmt.Sprintf(`{"query": {"constant_score": {"filter": {"bool": {"must":[{"match": {"title": "%s"}},{"match": {"authorsName": "%s"}},{"range": {"price": {"gte": %f, "lte": %f} }}]}}}}}`, title, author, priceStart, piceEnd)
	myJson := bytes.NewBuffer([]byte(s))

	req, err := http.NewRequest("GET", e.Url+"_search/", myJson)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	return e.Client.Do(req)
}

func (e *ElasticHandler) Store() (resp1 *http.Response, resp2 *http.Response, err error) {
	s1 := fmt.Sprintf(`{"query": {"match_all": {}}}`)
	myJson := bytes.NewBuffer([]byte(s1))

	req, err1 := http.NewRequest("GET", e.Url+"_count/", myJson)
	if err1 != nil {
		return nil, nil, err1
	}
	req.Header.Set("Content-Type", "application/json")
	resp1, err1 = e.Client.Do(req)
	if err1 != nil {
		return nil, nil, err1
	}

	s2 := fmt.Sprintf(`{"aggs" : {"authors_count" : {"cardinality" : {"field" : "authorsName.keyword"}}}}`)
	myJson2 := bytes.NewBuffer([]byte(s2))

	resp2, err2 := e.Client.Post(e.Url+"_search/", "application/json", myJson2)
	if err2 != nil {
		return nil, nil, err2
	}
	return resp1, resp2, nil
}
