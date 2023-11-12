package lru

import (
    "math/rand"
    "testing"
)

const (
    randSeed = 42
    length   = 10
)

func reverseArray(s []int) {
    for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
        s[i], s[j] = s[j], s[i]
    }
}

func TestLRU_size(t *testing.T) {
    l := NewLRU[int, int](0)
    if l.Len() != 0 {
        t.Error("size should be 0")
    }
    l = NewLRU[int, int](length)
    if l.Len() != 0 {
        t.Error("size should be 0")
    }
    if evicted := l.Set(1, 1); evicted != nil {
        t.Error("evicted should be nil")
    }
    if l.Len() != 1 {
        t.Error("size should be 1")
    }
}

func TestLRU_function(t *testing.T) {
    l := NewLRU[int, int](length)
    rand.Seed(randSeed)
    numbers := make([]int, length)
    for i := 0; i < length; i++ {
        numbers[i] = rand.Int()
        if evicted := l.Set(i, numbers[i]); evicted != nil {
            t.Fatalf("evicted should be nil")
        }
        if expected := l.Peek(i); *expected != numbers[i] {
            t.Fatalf("expected %v, got %v", numbers[i], expected)
        }
        if expected := l.Get(i); *expected != numbers[i] {
            t.Fatalf("expected %v, got %v", numbers[i], expected)
        }
    }
    reverseArray(numbers)
    entrys := l.GetOrderedEntrys()
    if len(entrys) != len(numbers) {
        t.Fatalf("Size is not equal expected %v, got %v", len(numbers), len(entrys))
    }
    for i := 0; i < length; i++ {
        if *entrys[i].Value != numbers[i] {
            t.Fatalf("Element is not equal, expected %v, got %v", numbers[i], entrys[i].Value)
        }
    }
}

func BenchmarkLRU_Rand(b *testing.B) {
    l := NewLRU[int64, int64](8192)
    rand.Seed(randSeed)
    trace := make([]int64, b.N*2)
    for i := 0; i < b.N*2; i++ {
        trace[i] = rand.Int63() % 32768
    }

    b.ResetTimer()

    var hit, miss int
    for i := 0; i < 2*b.N; i++ {
        if i%2 == 0 {
            l.Set(trace[i], trace[i])
        } else {
            if l.Get(trace[i]) == nil {
                miss++
            } else {
                hit++
            }
        }
    }
    b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}
