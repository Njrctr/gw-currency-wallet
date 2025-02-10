package cache

type Cache interface {
	Get(key string) (float64, bool)
	Set(key string, value float64)
}
