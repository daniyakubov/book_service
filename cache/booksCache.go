package cache

type Action struct {
	method string
	route  string
}

//the key is the username, and the value is it's actions
type BooksCache interface {
	Set(key string, value []Action)
	Get(key string) []Action
}
