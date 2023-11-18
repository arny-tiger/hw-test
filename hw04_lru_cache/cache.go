package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	keys     map[*ListItem]Key
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item)
		item.Value = value
		return true
	}
	if c.queue.Len() == c.capacity {
		tail := c.queue.Back()
		delete(c.items, c.keys[tail])
		delete(c.keys, tail)
		c.queue.Remove(tail)
	}
	newItem := c.queue.PushFront(value)
	c.items[key] = newItem
	c.keys[newItem] = key
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	val, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.MoveToFront(val)
	return val.Value, true
}

func (c *lruCache) Clear() {
	for k, item := range c.items {
		c.queue.Remove(item)
		delete(c.items, k)
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		keys:     make(map[*ListItem]Key, capacity),
	}
}
