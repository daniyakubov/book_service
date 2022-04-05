package cache

// Cache the key is the username, and the value is its last 3 actions
type Cache interface {
	Push(key string, value string) error
	Get(key string) ([]string, error)
}
