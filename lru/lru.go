package lru

import "container/list"

// Cache is an LRU cache
type Cache struct {
	// maxEntries is the maximum number of cache entries
	maxEntries int

	// linkedList stores a doubly linked list
	linkedList *list.List

	// cache stores all cache data
	cacheByKey map[string]*list.Element
}

// entry is the value of the linked list defined
type entry struct {
	key string
	value interface{}
}

// New creates a new cache.
func New(maxEntries int) *Cache {
	return &Cache{
		maxEntries: maxEntries,
		linkedList: list.New(),
		cacheByKey:      make(map[string]*list.Element),
	}
}

// Set sets a new entry to the cache.
func (c *Cache) Set(key string, value interface{})  {
	if el, ok := c.cacheByKey[key]; ok {
		c.linkedList.MoveToFront(el)
		el.Value.(*entry).value = value
		return
	}

	el := c.linkedList.PushFront(&entry{key: key, value: value})
	c.cacheByKey[key] = el

	if c.linkedList.Len() > c.maxEntries {
		c.cleanUp()
	}
}

// Get retries cache's value by key.
func (c *Cache) Get(key string) (value interface{}, ok bool) {
	if el, hit := c.cacheByKey[key]; hit {
		c.linkedList.MoveToFront(el)

		value = el.Value.(*entry).value
		ok = true
		return
	}

	return
}

// GetAllKeys returns a list all keys available in cache.
func (c *Cache) GetAllKeys() []string {
	keys := make([]string, 0, len(c.cacheByKey))

	for k := range c.cacheByKey {
		keys = append(keys, k)
	}

	return keys
}

// Delete deletes a cache entry by key.
func (c *Cache) Delete(key string) {
	if el, hit := c.cacheByKey[key]; hit {
		c.delete(el)
	}
}

func (c *Cache) delete(el *list.Element) {
	c.linkedList.Remove(el)

	entry := el.Value.(*entry)
	delete(c.cacheByKey, entry.key)
}

func (c *Cache) cleanUp() {
	c.delete(c.linkedList.Back())
}
