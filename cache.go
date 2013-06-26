package cache

import (
	"sync"
)

type Cache interface {
	Get(string) string
	Set(string, string)
}

type concurrentCache struct {
	buckets  *[BUCKET_SIZE]Cache
	hashFunc func(str string) int
}

func (c *concurrentCache) Get(key string) string {
	idx := indexFor(c.hashFunc(key))
	return c.buckets[idx].Get(key)
}

func (c *concurrentCache) Set(key, val string) {
	idx := indexFor(c.hashFunc(key))
	c.buckets[idx].Set(key, val)
}

type standardCache struct {
	m map[string]string
	l sync.RWMutex
}

func (c *standardCache) Get(key string) string {
	c.l.RLock()
	defer c.l.RUnlock()

	return c.m[key]
}

func (c *standardCache) Set(key, val string) {
	c.l.Lock()
	defer c.l.Unlock()

	c.m[key] = val
}

func NewConcurrentCache() Cache {
	var b [BUCKET_SIZE]Cache

	for i := 0; i < BUCKET_SIZE; i++ {
		b[i] = NewCache()
	}

	return &concurrentCache{
		buckets: &b,
		hashFunc: func(str string) int {
			return hash(str, 0)
		},
	}
}

func NewCache() Cache {
	m := make(map[string]string)
	return &standardCache{m: m}
}

// Lets hard code this for now.
// No idea of the implications a for the algorithm if this can vary...
const BUCKET_SIZE = 16

func hash(str string, idx int) int {
	return hashInt(int(str[idx]))
}

func hashInt(h int) int {
	h ^= (h >> 20) ^ (h >> 12)
	return h ^ (h >> 7) ^ (h >> 4)
}

func indexFor(h int) int {
	return h & (BUCKET_SIZE - 1)
}
