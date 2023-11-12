package fifo

import (
    "math/rand"
    "testing"
)

const (
    randSeed = 42
    length   = 10
)

func reverseSlice(s []int) {
    for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
        s[i], s[j] = s[j], s[i]
    }
}

func TestFIFO(t *testing.T) {
    f := NewFIFO[int, int](0)

    if e := f.Set(0, 1); e == nil || e.Value != 1 {
        t.Fatalf("firset element should be evicted")
    }

    if e := f.Get(1); e != nil {
        t.Fatalf("fifo must be empty")
    }

    if 0 != f.Len() {
        t.Fatalf("fifo size must be zero")
    }

    f = NewFIFO[int, int](length)
    rand.Seed(randSeed)
    numbers := make([]int, length)
    for i := 0; i < length; i++ {
        numbers[i] = rand.Int()
        evicted := f.Set(i, numbers[i])
        if evicted != nil {
            t.Fatalf("evicted element must be nil")
        }
    }
    if f.Len() != length {
        t.Fatalf("fifo size must be %d", length)
    }
    for i := 0; i < length; i++ {
        if value := f.Get(i); value == nil || *value != numbers[i] {
            t.Fatalf("element %d must be %d", i, numbers[i])
        }
    }
}

func TestFIFO_Order(t *testing.T) {
    f := NewFIFO[int, int](length)
    numbers := make([]int, length)
    for i := 0; i < length; i++ {
        numbers[i] = rand.Int()
        evicted := f.Set(i, numbers[i])
        if evicted != nil {
            t.Fatalf("evicted element must be nil")
        }
    }
    valuesOrdered := f.GetOrderedValues()
    if len(valuesOrdered) != len(numbers) {
        t.Fatalf("fifo size must be %d", length)
    }
    reverseSlice(valuesOrdered)
    for i := 0; i < length; i++ {
        if valuesOrdered[i] != numbers[i] {
            t.Fatalf("element %d must be %d, but get %d", i, numbers[i], valuesOrdered[i])
        }
    }
}

func TestFIFO_function(t *testing.T) {
    f := NewFIFO[int, int](2)
    if f.Len() != 0 {
        t.Fatalf("fifo size must be %d", 0)
    }
    if f.Set(1, 1) != nil {
        t.Fatalf("evicted element must be nil")
    }
    if f.Set(2, 2) != nil {
        t.Fatalf("evicted element must be nil")
    }
    if f.Len() != 2 {
        t.Fatalf("fifo size must be %d", 2)
    }
    if evicted := f.Set(3, 3); evicted == nil || evicted.Key != 1 || evicted.Value != 1 {
        t.Fatalf("evicted element must be {1, 1}")
    }
    if evicted := f.Set(4, 4); evicted == nil || evicted.Key != 2 || evicted.Value != 2 {
        t.Fatalf("evicted element must be {2, 2}")
    }
}
