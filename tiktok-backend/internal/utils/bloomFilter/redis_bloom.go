package bloomFilter

import (
    "context"
    "hash/fnv"

    "github.com/go-redis/redis/v8"
)

// BloomFilter struct
type RedisBloomFilter struct {
    redisClient *redis.Client
    key         string
    hashes      int
    size        int
}

// NewBloomFilter creates a new Bloom Filter
func NewBloomFilter(redisClient *redis.Client, key string, hashes int) *RedisBloomFilter {
    return &RedisBloomFilter{
        redisClient: redisClient,
        key:         key,
        hashes:      hashes,
    }
}

// hash calculates a hash for a given string and a seed
func hash(s string, seed uint32) uint32 {
    h := fnv.New32a()
    h.Write([]byte(s))
    return h.Sum32() + seed
}

// Add adds an element to the Bloom Filter
func (bf *RedisBloomFilter) Add(ctx context.Context, element string) error {
    for i := 0; i < bf.hashes; i++ {
        position := hash(element, uint32(i)) % uint32(bf.size)
        err := bf.redisClient.SetBit(ctx, bf.key, int64(position), 1).Err()
        if err != nil {
            return err
        }
    }
    return nil
}

// Check checks if an element might be in the set
func (bf *RedisBloomFilter) Check(ctx context.Context, element string) (bool, error) {
    for i := 0; i < bf.hashes; i++ {
        position := hash(element, uint32(i)) % uint32(bf.size)
        bit, err := bf.redisClient.GetBit(ctx, bf.key, int64(position)).Result()
        if err != nil {
            return false, err
        }
        if bit == 0 {
            return false, nil
        }
    }
    return true, nil
}

func (bf *RedisBloomFilter) Clear(ctx context.Context) error {
    _, err := bf.redisClient.Del(ctx, bf.key).Result()
    return err
}
