package syncmap

import (
    "math/rand"
    "sync"
    "testing"
)

const (
    randSeed = 42
)

var once sync.Once

func nrands(n int) []int {
    once.Do(func() {
        rand.Seed(randSeed)
    })
    nums := make([]int, n)
    for i := 0; i < n; i++ {
        nums[i] = rand.Int()
    }
    return nums
}

type RWMap[K comparable, V any] struct {
    m        sync.RWMutex
    internal map[K]V
}

func NewRWMap[K comparable, V any]() *RWMap[K, V] {
    return &RWMap[K, V]{
        internal: make(map[K]V),
    }
}

func (rm *RWMap[K, V]) Get(key K) (value V, ok bool) {
    rm.m.RLock()
    defer rm.m.RUnlock()
    value, ok = rm.internal[key]
    return
}

func (rm *RWMap[K, V]) Set(key K, value V) {
    rm.m.Lock()
    defer rm.m.Unlock()
    rm.internal[key] = value
}

func (rm *RWMap[K, V]) Delete(key K) {
    rm.m.Lock()
    defer rm.m.Unlock()
    delete(rm.internal, key)
}

type SyncMap struct {
    internal sync.Map
}

func NewSyncMap() *SyncMap {
    return &SyncMap{}
}

func (sm *SyncMap) Get(key any) (value any, ok bool) {
    value, ok = sm.internal.Load(key)
    return
}

func (sm *SyncMap) Set(key any, value any) {
    sm.internal.Store(key, value)
}

func (sm *SyncMap) Delete(key any) {
    sm.internal.Delete(key)
}

func BenchmarkRegular(b *testing.B) {
    nums := nrands(b.N)
    rm := NewRWMap[int, int]()
    for idx, v := range nums {
        rm.Set(idx, v)
    }
    b.ResetTimer()
    for _, v := range nums {
        rm.Delete(v)
    }
}

func BenchmarkSync(b *testing.B) {
    nums := nrands(b.N)
    sm := NewSyncMap()
    for idx, v := range nums {
        sm.Set(idx, v)
    }
    b.ResetTimer()
    for _, v := range nums {
        sm.Delete(v)
    }
}
