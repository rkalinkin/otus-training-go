package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value any) bool
	Get(key Key) (any, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mutex    sync.RWMutex
}

func (c *lruCache) Set(key Key, value any) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item, exists := c.items[key]; exists {
		entry := item.Value.(*cacheEntry)
		entry.value = value
		c.queue.MoveToFront(item)
		return true
	}

	entry := &cacheEntry{
		key:   key,
		value: value,
	}
	item := c.queue.PushFront(entry)
	c.items[key] = item

	if c.queue.Len() > c.capacity {
		back := c.queue.Back()
		backEntry := back.Value.(*cacheEntry)
		delete(c.items, backEntry.key)
		c.queue.Remove(back)
	}

	return false
}

func (c *lruCache) Get(key Key) (any, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item, exists := c.items[key]; exists {
		c.queue.MoveToFront(item)
		entry := item.Value.(*cacheEntry)
		return entry.value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

type cacheEntry struct {
	key   Key
	value any
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
