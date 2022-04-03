package cache

type Actions struct {
	method1 string
	route1  string
	method2 string
	route2  string
	method3 string
	route3  string
}

//the key is the username, and the value is it's last 3 actions
type BooksCache interface {
	Set(key string, value Actions)
	Get(key string) *Actions
}
