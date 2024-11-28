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
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	elem, ok := c.items[key]
	if ok {
		if len(c.items) == c.capacity {
			old := c.queue.Back()
			c.queue.Remove(old)
			str, ok := old.Key.(Key)
			if ok {
				delete(c.items, str)
			}
		}
		elem.Value = value
		elem.Key = key
		c.queue.MoveToFront(elem)
		return true
	}
	el := c.queue.PushFront(value)
	c.items[key] = el
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	elem, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.MoveToFront(elem)
	return elem.Value, true

}

func (c *lruCache) Clear() {

}
