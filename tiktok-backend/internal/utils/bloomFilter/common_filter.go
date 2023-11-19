package bloomFilter

import (
    "context"

    "github.com/bits-and-blooms/bloom/v3"
)

// CommonBloomFilter construct the Bloom filter capable of receiving 1 million elements with a false-positive
// rate of 1% in the following manner.
// filter := bloom.NewWithEstimates(1000000, 0.01)
// New creates a new Bloom filter with _m_ bits and _k_ hashing functions
// We force _m_ and _k_ to be at least one to avoid panics.
// f := New(1000, 4)
// m bits, k hashes for a set of size n for the actual false positive rate(fp)
// m, k := bloom.EstimateParameters(n, fp)
// fpRate := bloom.EstimateFalsePositiveRate(m, k, n)
type CommonBloomFilter struct {
    filter *bloom.BloomFilter
    key    string
    hashes int
    size   int
}

func NewCommonBloomFilter(key string, hashes int, size int) *CommonBloomFilter {
    return &CommonBloomFilter{
        filter: bloom.New(uint(size), uint(hashes)),
        key:    key,
        hashes: hashes,
        size:   size,
    }
}

func (bf *CommonBloomFilter) Add(ctx context.Context, element string) error {
    bf.filter.Add([]byte(element))
    return nil
}

func (bf *CommonBloomFilter) Check(ctx context.Context, element string) (bool, error) {
    return bf.filter.Test([]byte(element)), nil
}

func (bf *CommonBloomFilter) Clear(ctx context.Context) error {
    bf.filter.ClearAll()
    return nil
}
