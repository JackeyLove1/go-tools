package fifo

import (
    "sync"

    "go-tools/tools/cache_v1/list"
)

// FIFO TODO: use bucket to improve performance
type FIFO[K comparable, V any] struct {
    m     sync.Mutex
    ll    *list.List[*entry[K, V]]
    cache map[K]*list.Element[*entry[K, V]]
    size  int
}

type entry[K comparable, V any] struct {
    Key   K
    Value *V
}

type Evicted[K comparable, V any] struct {
    Key   K
    Value V
}

func (L *FIFO[K, V]) Get(key K) *V {
    L.m.Lock()
    defer L.m.Unlock()
    if e, ok := L.cache[key]; ok {
        return e.Value.Value
    }
    return nil
}

func (L *FIFO[K, V]) Set(key K, value V) *Evicted[K, V] {
    L.m.Lock()
    defer L.m.Unlock()

    if L.size < 1 {
        return &Evicted[K, V]{key, value}
    }

    if e, ok := L.cache[key]; ok {
        e.Value.Value = &value
        return nil
    }
    e := L.ll.Back()
    targetEntry := e.Value
    evictedKey := targetEntry.Key
    evictedValue := targetEntry.Value
    if evictedValue != nil {
        delete(L.cache, evictedKey)
    }
    targetEntry.Key = key
    targetEntry.Value = &value
    L.ll.MoveToFront(e)
    L.cache[key] = e
    // println("key:", key, " value:", *e.Value.Value, " cache size:", len(L.cache))
    if evictedValue != nil {
        return &Evicted[K, V]{Key: evictedKey, Value: *evictedValue}
    }
    return nil
}

func (L *FIFO[K, V]) Len() int {
    L.m.Lock()
    defer L.m.Unlock()
    return len(L.cache)
}

func (L *FIFO[K, V]) Remove(key K) *V {
    L.m.Lock()
    defer L.m.Unlock()

    if e, ok := L.cache[key]; ok {
        value := e.Value.Value
        L.ll.MoveToBack(e)
        delete(L.cache, key)
        return value
    }
    return nil
}

func NewFIFO[K comparable, V any](size int) *FIFO[K, V] {
    fifo := &FIFO[K, V]{
        ll:    list.NewList[*entry[K, V]](),
        cache: make(map[K]*list.Element[*entry[K, V]]),
        size:  size,
    }
    for i := 0; i < size; i++ {
        fifo.ll.PushBack(&entry[K, V]{})
    }
    return fifo
}

func (L *FIFO[K, V]) GetCacheKeys() []K {
    L.m.Lock()
    defer L.m.Unlock()
    keys := make([]K, 0, len(L.cache))
    for k, _ := range L.cache {
        keys = append(keys, k)
    }
    return keys
}

func (L *FIFO[K, V]) GetCacheValues() []V {
    L.m.Lock()
    defer L.m.Unlock()
    values := make([]V, 0, len(L.cache))
    for _, v := range L.cache {
        values = append(values, *v.Value.Value)
    }
    return values
}

func (L *FIFO[K, V]) GetEntrys() []*entry[K, V] {
    L.m.Lock()
    defer L.m.Unlock()
    entries := make([]*entry[K, V], 0, len(L.cache))
    for e := L.ll.Front(); e != nil; e = e.Next() {
        if _, ok := L.cache[e.Value.Key]; ok {
            entries = append(entries, e.Value)
        }
    }
    return entries
}

func (L *FIFO[K, V]) GetOrderedKeys() []K {
    entrys := L.GetEntrys()
    keys := make([]K, 0, len(entrys))
    for _, e := range entrys {
        keys = append(keys, e.Key)
    }
    return keys
}

func (L *FIFO[K, V]) GetOrderedValues() []V {
    entrys := L.GetEntrys()
    values := make([]V, 0, len(entrys))
    for _, e := range entrys {
        values = append(values, *e.Value)
    }
    return values
}
