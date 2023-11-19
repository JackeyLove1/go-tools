package bloomFilter

import "context"

type Filter interface {
    Add(ctx context.Context, key string) error
    Check(ctx context.Context, key string) (bool, error)
    Clear(ctx context.Context) error
}
