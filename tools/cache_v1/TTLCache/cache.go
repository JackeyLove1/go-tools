package TTLCache

import (
    "sync"
    "time"
)

type Cache struct {
    items sync.Map // key: any, value: *item
    close chan struct{}
}

type item struct {
    data    any
    expires int64
}

func NewCache(cleaningInterval time.Duration) *Cache {
    cache := &Cache{
        close: make(chan struct{}),
    }
    go func() {
        ticker := time.NewTicker(cleaningInterval)
        defer ticker.Stop()
        for {
            select {
            case <-cache.close:
                return
            case <-ticker.C:
                now := time.Now().UnixNano()
                cache.items.Range(func(key, value any) bool {
                    item := value.(*item)
                    if item.expires > 0 && item.expires < now {
                        cache.items.Delete(key)
                    }
                    return true
                })
            }
        }
    }()
    return cache
}

func (c *Cache) Set(key any, value any, duration time.Duration) {
    var expires int64
    if duration > 0 {
        expires = time.Now().Add(duration).UnixNano()
    }
    c.items.Store(key, &item{data: value, expires: expires})
}

func (c *Cache) Get(key any) (any, bool) {
    value, ok := c.items.Load(key)
    if !ok {
        return nil, false
    }
    item := value.(*item)
    if item.expires > 0 && item.expires < time.Now().UnixNano() {
        c.items.Delete(key)
        return nil, false
    }
    return item.data, true
}

func (c *Cache) Delete(key any) {
    c.items.Delete(key)
}

func (c *Cache) Close() {
    c.close <- struct{}{}
    c.items = sync.Map{}
}

func (c *Cache) Range(f func(key, value any) bool) {
    now := time.Now().UnixNano()
    fn := func(key, value any) bool {
        item := value.(*item)
        if item.expires > 0 && item.expires < now {
            return false
        }
        return f(key, item.data)
    }
    c.items.Range(fn)
}
