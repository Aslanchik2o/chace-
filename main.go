package cache

import (
	"sync"
	"time"
)

// Cache - структура, хранящая элементы в карте
type Cache struct {
	items sync.Map
}

// Item - структура, хранящая значение и время истечения
type Item struct {
	Value      interface{}
	Expiration time.Time
}

// New - функция, создающая новый экземпляр Cache
func New() *Cache {
	return &Cache{}
}

// Set - метод, добавляющий элемент в кэш по ключу с заданным TTL
func (c *Cache) Set(key interface{}, value interface{}, ttl time.Duration) {
	item := Item{
		Value:      value,
		Expiration: time.Now().Add(ttl),
	}
	c.items.Store(key, item)
	// запускаем функцию, которая удалит элемент из кэша по истечении TTL
	time.AfterFunc(ttl, func() {
		c.Delete(key)
	})
}

// Get - метод, возвращающий элемент из кэша по ключу
func (c *Cache) Get(key interface{}) (interface{}, bool) {
	item, ok := c.items.Load(key)
	if !ok {
		return nil, false
	}
	// проверяем, не истек ли TTL элемента
	if item.(Item).Expiration.Before(time.Now()) {
		return nil, false
	}
	return item.(Item).Value, true
}

// Delete - метод, удаляющий элемент из кэша по ключу
func (c *Cache) Delete(key interface{}) {
	c.items.Delete(key)
}
