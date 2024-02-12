package cache

type CacheInterface interface {
	Set(key, value string) error
	Get(key string) (string, error)
}
