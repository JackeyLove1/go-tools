package twoqueue

import (
    "go-tools/tools/cache_v1/fifo"
    "go-tools/tools/cache_v1/lru"
)

const (
    Default2QRecentRatio  = 0.25
    Default2QGhostEntries = 0.50
)

type TwoQueue[K comparable, V any] struct {
    recent        *fifo.FIFO[K, V]
    recentEvicted *fifo.FIFO[K, struct{}]
    frequent      *lru.LRU[K, V]
}

type Evicted[K comparable, V any] struct {
    Key   K
    Value V
}

func fromLruEvicted[K comparable, V any](e *lru.Evicted[K, V]) *Evicted[K, V] {
    if e == nil {
        return nil
    }
    return &Evicted[K, V]{e.Key, e.Value}
}

func fromFifoEvicted[K comparable, V any](e *fifo.Evicted[K, V]) *Evicted[K, V] {
    if e == nil {
        return nil
    }
    return &Evicted[K, V]{e.Key, e.Value}
}

func (L *TwoQueue[K, V]) Get(key K) *V {
    if e := L.frequent.Get(key); e != nil {
        return e
    }
    return L.recent.Get(key)
}

func (L *TwoQueue[K, V]) Set(key K, value V) *Evicted[K, V] {
    if e := L.frequent.Peek(key); e != nil {
        return fromLruEvicted(L.frequent.Set(key, value))
    }
    if L.recentEvicted.Get(key) != nil {
        L.recentEvicted.Remove(key)
        return fromLruEvicted(L.frequent.Set(key, value))
    }
    if re := L.recent.Set(key, value); re != nil {
        L.recent.Set(key, value)
        return fromFifoEvicted(re)
    }
    return nil
}

func (L *TwoQueue[K, V]) Len() int {
    return L.recent.Len() + L.frequent.Len()
}

func (L *TwoQueue[K, V]) Peek(key K) *V {
    if e := L.frequent.Peek(key); e != nil {
        return e
    }
    return L.recent.Get(key)
}

func (L *TwoQueue[K, V]) Remove(key K) *V {
    if e := L.frequent.Remove(key); e != nil {
        return e
    }
    return L.recent.Remove(key)
}

func NewParams[K comparable, V any](Kin int, Kout int, size int) *TwoQueue[K, V] {
    return &TwoQueue[K, V]{
        recent:        fifo.NewFIFO[K, V](Kin),
        recentEvicted: fifo.NewFIFO[K, struct{}](Kout),
        frequent:      lru.NewLRU[K, V](size),
    }
}

func New[K comparable, V any](size int) *TwoQueue[K, V] {
    return NewParams[K, V](
        int(Default2QRecentRatio*float64(size)),
        int(Default2QGhostEntries*float64(size)),
        int((1-Default2QRecentRatio)*float64(size)),
    )
}
