package models

type Request struct {
	Data  *Hit
	Route string
}

type GetBookRequest struct {
	Data  *Hit
	Route string
}

func NewRequest(data *Hit, route string) Request {
	return Request{
		Data:  data,
		Route: route,
	}
}
