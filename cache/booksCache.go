package cache

type Action struct {
	method string
	route  string
}

type CacheBooks interface {
	Set(key string, value *Action)
	Get(key string) *Action
}
