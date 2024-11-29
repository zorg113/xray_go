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

// создание кэша
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// Реализация интерфейса добавления элемента в кэш
func (c *lruCache) Set(key Key, value interface{}) bool {
	elem, ok := c.items[key]
	if ok {
		elem.Value = value
		elem.Key = key
		c.queue.MoveToFront(elem)
		return true
	}
	if len(c.items) == c.capacity {
		old := c.queue.Back()
		c.queue.Remove(old)
		str, ok := old.Key.(Key)
		if ok {
			delete(c.items, str)
		}
	}
	el := c.queue.PushFront(value)
	el.Key = key
	c.items[key] = el
	return false
}

// Извлечение элемента из кэша
func (c *lruCache) Get(key Key) (interface{}, bool) {
	elem, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.MoveToFront(elem)
	return elem.Value, true

}

// очистка кэша
func (c *lruCache) Clear() {

}
