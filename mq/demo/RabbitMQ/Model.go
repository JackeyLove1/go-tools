package RabbitMQ

import (
    "fmt"
    "math/rand"
    "sync"
    "sync/atomic"
    "time"
)

type RandomNumber struct {
    RandNum int
    Idx     uint32
}

func (r *RandomNumber) String() string {
    return fmt.Sprintf("Idx: %d, RandNum: %d", r.Idx, r.RandNum)
}

var random *rand.Rand
var once sync.Once
var count atomic.Uint32

func init() {
    once.Do(func() {
        count.Store(0)
    })
    random = rand.New(rand.NewSource(time.Now().Unix()))
}

func GenerateRandomNumber() RandomNumber {
    defer count.Add(1)
    return RandomNumber{
        RandNum: random.Int(),
        Idx:     count.Load(),
    }
}
