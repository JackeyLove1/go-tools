package memtable

import (
    "sync"

    "go-tools/db/LSM/utils"
)

//TODO: use key value as pointer to accelerate
type MemTable struct {
    KV     *SkipList[string, utils.Item]
    RWLock sync.RWMutex
}

func NewMemTable() *MemTable {
    return &MemTable{
        KV: New[string, utils.Item](String),
    }
}

func (m *MemTable) GetCount() int {
    m.RWLock.RLock()
    defer m.RWLock.RUnlock()
    return m.KV.Len()
}

func (m *MemTable) Get(key string) (utils.Item, utils.SearchResult) {
    m.RWLock.RLock()
    defer m.RWLock.RUnlock()

    item, ok := m.KV.GetValue(key)
    if !ok || item.Key != key {
        return utils.Item{}, utils.None
    }
    if item.Deleted == true {
        return utils.Item{}, utils.Deleted
    }
    return item, utils.Found
}

func (m *MemTable) Set(key string, value utils.Item) {
    m.RWLock.Lock()
    defer m.RWLock.Unlock()
    m.KV.Set(key, value)
}

func (m *MemTable) Delete(key string) {
    m.Set(key, utils.Item{Key: key, Value: []byte(nil), Deleted: true})
}

func (m *MemTable) GetItems() []utils.Item {
    m.RWLock.RLock()
    defer m.RWLock.RUnlock()
    items := make([]utils.Item, 0, m.KV.Len())
    for iter := m.KV.Front(); iter != nil; iter = iter.Next() {
        items = append(items, iter.Value)
    }
    return items
}

func (m *MemTable) GetValues() [][]byte {
    items := m.GetItems()
    values := make([][]byte, 0, len(items))
    for _, item := range items {
        values = append(values, item.Value)
    }
    return values
}

func (m *MemTable) GetKeys() []string {
    items := m.GetItems()
    keys := make([]string, 0, len(items))
    for _, item := range items {
        keys = append(keys, item.Key)
    }
    return keys
}
