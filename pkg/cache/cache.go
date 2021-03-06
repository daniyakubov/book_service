package cache

// Cache the key is the username, and the value is its last 3 actions
type Cache interface {
	Push(key string, value string)
	Get(key string) string
}
