package lru

import (
    "sync"

    "go-tools/tools/cache_v1/list"
)

type entry[K comparable, V any] struct {
    Key   K
    Value *V
}

type Evicted[K comparable, V any] struct {
    Key   K
    Value V
}

type LRU[K comparable, V any] struct {
    m     sync.Mutex
    cache map[K]*list.Element[*entry[K, V]]
    ll    *list.List[*entry[K, V]]
    size  int
}

func NewLRU[K comparable, V any](size int) *LRU[K, V] {
    lru := &LRU[K, V]{
        cache: make(map[K]*list.Element[*entry[K, V]]),
        ll:    list.NewList[*entry[K, V]](),
        size:  size,
    }
    for i := 0; i < size; i++ {
        lru.ll.PushBack(&entry[K, V]{})
    }
    return lru
}

func (l *LRU[K, V]) Get(key K) *V {
    l.m.Lock()
    defer l.m.Unlock()
    if e, ok := l.cache[key]; ok {
        return e.Value.Value
    }
    return nil
}

func (l *LRU[K, V]) Set(key K, value V) *Evicted[K, V] {
    if l.size < 1 {
        return &Evicted[K, V]{key, value}
    }
    l.m.Lock()
    defer l.m.Unlock()
    if e, ok := l.cache[key]; ok {
        preValue := e.Value.Value
        e.Value.Value = &value
        l.ll.MoveToFront(e)
        return &Evicted[K, V]{key, *preValue}
    }
    e := l.ll.Back()
    preValue := e.Value.Value
    if preValue != nil {
        delete(l.cache, e.Value.Key)
    }
    e.Value.Key = key
    e.Value.Value = &value
    l.ll.MoveToFront(e)
    l.cache[key] = e
    if preValue != nil {
        return &Evicted[K, V]{Key: e.Value.Key, Value: *preValue}
    }
    return nil
}

func (l *LRU[K, V]) Len() int {
    l.m.Lock()
    defer l.m.Unlock()
    return len(l.cache)
}

func (l *LRU[K, V]) Remove(key K) *V {
    l.m.Lock()
    defer l.m.Unlock()
    if e, ok := l.cache[key]; ok {
        preValue := e.Value.Value
        delete(l.cache, key)
        l.ll.Remove(e)
        return preValue
    }
    return nil
}

func (l *LRU[K, V]) Peek(key K) *V {
    l.m.Lock()
    defer l.m.Unlock()
    if e, ok := l.cache[key]; ok {
        return e.Value.Value
    }
    return nil
}

func (l *LRU[K, V]) GetOrderedEntrys() []*entry[K, V] {
    l.m.Lock()
    defer l.m.Unlock()
    entrys := make([]*entry[K, V], 0, len(l.cache))
    for e := l.ll.Front(); e != nil; e = e.Next() {
        key := e.Value.Key
        if _, ok := l.cache[key]; ok {
            entrys = append(entrys, e.Value)
        }
    }
    return entrys
}
