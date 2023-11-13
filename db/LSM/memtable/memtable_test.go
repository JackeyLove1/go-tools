package memtable

import (
    "math/rand"
    "strconv"
    "sync"
    "testing"

    "go-tools/db/LSM/utils"
)

const length = 10000
const randSeed = 42

var once sync.Once

func GenerateRandomString(length int) string {
    once.Do(func() {
        rand.Seed(randSeed)
    })
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    randomBytes := make([]byte, length)
    charsetLength := len(charset)
    for i := 0; i < length; i++ {
        randomIndex := rand.Int() % charsetLength
        randomBytes[i] = charset[randomIndex]
    }
    return string(randomBytes)
}

func TestMemTableBasic(t *testing.T) {
    m := NewMemTable()
    // Test Set Get
    for i := 0; i < length; i++ {
        s := strconv.Itoa(i)
        m.Set(s, utils.Item{Key: s, Value: []byte(s), Deleted: false})
        if item, result := m.Get(s); result != utils.Found || string(item.Value) != s {
            t.FailNow()
        }
    }
    // Test Delete
    for i := 0; i < length; i++ {
        s := strconv.Itoa(i)
        m.Delete(s)
        if item, result := m.Get(s); result != utils.Deleted || len(item.Value) != 0 {
            t.FailNow()
        }
    }
    // Test None
    for i := length; i < length*2; i++ {
        s := strconv.Itoa(i)
        if item, result := m.Get(s); result != utils.None || len(item.Value) != 0 {
            t.FailNow()
        }
    }

    for i := 0; i < length; i++ {
        s := strconv.Itoa(i)
        m.Set(s, utils.Item{Key: s, Value: []byte(s), Deleted: false})
    }
    // Test Iterator
    keys := m.GetKeys()
    values := m.GetValues()
    if len(keys) != length || len(values) != length {
        t.Fatalf("keys and values length is not equal to length")
    }
    kvs := make(map[string]string, length)
    for i := 0; i < length; i++ {
        kvs[keys[i]] = string(values[i])
    }
    for i := 0; i < length; i++ {
        k := strconv.Itoa(i)
        if v, ok := kvs[k]; !ok || v != k {
            t.FailNow()
        }
    }
}

func BenchmarkBasic(b *testing.B) {
    m := NewMemTable()
    for i := 0; i < b.N; i++ {
        s := strconv.Itoa(i)
        m.Set(s, utils.Item{Key: s, Value: []byte(s), Deleted: false})
        if item, result := m.Get(s); result != utils.Found || string(item.Value) != s {
            b.FailNow()
        }
    }
}
