package TTLCache

import (
    "math/rand"
    "sync/atomic"
    "testing"
    "time"
)

const (
    intervalTime = 100
    interval     = intervalTime * time.Millisecond
    length       = 100
    Seed         = 1e9 + 7
    Size         = 1e6
)

func GetRandomNumbers(size int) []int {
    numbers := make([]int, size)
    for i := 0; i < size; i++ {
        numbers[i] = rand.Int() % Seed
    }
    return numbers
}

func TestGetSet(t *testing.T) {
    c := NewCache(interval)
    defer c.Close()
    c.Set("sticky", "forever", 0)
    c.Set("how", "are you", interval/2)
    if how, ok := c.Get("how"); !ok || how != "are you" {
        t.Fatalf("Key how not found")
    }
    time.Sleep(interval)
    if sticky, ok := c.Get("sticky"); !ok || sticky != "forever" {
        t.Fatalf("Key sticky not found")
    }
    if _, ok := c.Get("how"); ok {
        t.Fatalf("Key how should be cleaned for timeout")
    }
}

func TestDelete(t *testing.T) {
    c := NewCache(interval)
    defer c.Close()
    c.Set("sticky", "forever", 0)
    c.Set("how", "are you", interval/2)
    c.Delete("how")
    c.Delete("sticky")
    if _, ok := c.Get("how"); ok {
        t.Fatalf("Key how should be cleaned for timeout")
    }
    if _, ok := c.Get("sticky"); ok {
        t.Fatalf("Key sticky should be cleaned for timeout")
    }
}

func TestRange(t *testing.T) {
    c := NewCache(interval)
    defer c.Close()
    for i := 0; i < length; i++ {
        number := rand.Int()
        c.Set(i, number, 0)
        if value, ok := c.Get(i); !ok || value != number {
            t.Fatalf("Key %d not found", i)
        }
    }
    var count int32 = 0
    c.Range(func(key, value any) bool {
        atomic.AddInt32(&count, 1)
        return true
    })
    if count != length {
        t.Fatalf("Count %d != %d", count, length)
    }
}

func BenchmarkSetGet1(b *testing.B) {
    c := NewCache(interval)
    defer c.Close()

    // keys := GetRandomNumbers(b.N)
    values := GetRandomNumbers(b.N)
    for i := 0; i < b.N; i++ {
        c.Set(i, values[i], 0)
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        value, ok := c.Get(i)
        if !ok || value != values[i] {
            b.FailNow()
        }
    }
}

func BenchmarkSetGet(b *testing.B) {
    c := NewCache(interval)
    defer c.Close()

    values := GetRandomNumbers(b.N)
    for i := 0; i < b.N; i++ {
        c.Set(i, values[i], 0)
    }
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            for i := 0; i < length; i++ {
                idx := rand.Int() % b.N
                c.Get(idx)
                c.Set(idx, rand.Int(), 0)
            }
        }
    })
}
