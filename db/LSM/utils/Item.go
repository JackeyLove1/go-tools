package utils

import "encoding/json"

type SearchResult int

const (
    None SearchResult = iota
    Deleted
    Found
)

type Item struct {
    Key     string
    Value   []byte
    Deleted bool
}

func GetItemValue[T any](item *Item) (any, error) {
    var value T
    err := json.Unmarshal(item.Value, &value)
    return value, err
}

func Convert[T any](value T) ([]byte, error) {
    return json.Marshal(value)
}

func EncodeItem(item Item) ([]byte, error) {
    return json.Marshal(item)
}

func DecodeItem(data []byte) (Item, error) {
    var item Item
    err := json.Unmarshal(data, &item)
    return item, err
}
