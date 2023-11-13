package wal

import (
    "os"
    "sync"
)

type Wal struct {
    f    *os.File
    path string
    lock sync.Locker
}
