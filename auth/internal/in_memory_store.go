package internal

import "sync"

// inmemory cache  store to keep track of the request cycle with authorization code
type cache struct {
	kvmap map[string]any
	mx    *sync.RWMutex
}

func (c *cache) Get(k string) any {
	if c == nil {
		return nil
	}
	c.mx.RLock()
	defer c.mx.RUnlock()
	return c.kvmap[k]
}

func (c *cache) Set(k string, v any) {
	if c == nil {
		return
	}
	c.mx.Lock()
	defer c.mx.Unlock()
	c.kvmap[k] = v
}

func (c *cache) Delete(k string) {
	if c == nil {
		return
	}
	c.mx.Lock()
	defer c.mx.Unlock()
	delete(c.kvmap, k)
}

var (
	cacheClient *cache = &cache{
		kvmap: make(map[string]any),
		mx:    &sync.RWMutex{},
	}
)
