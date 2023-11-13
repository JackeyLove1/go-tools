package sstable

import "os"

/*
index from the data block
0 ─────────────────────────────────────────────────────────►
◄───────────────────────────
          dataLen          ◄──────────────────
                                indexLen     ◄──────────────┐
┌──────────────────────────┬─────────────────┬──────────────┤
│                          │                 │              │
│          数据区           │   稀疏索引区     │    元数据     │
│                          │                 │              │
└──────────────────────────┴─────────────────┴──────────────┘
*/

// meta data in the back of data block
type MetaInfo struct {
    version     uint64
    dataStart   uint64
    dataLength  uint64
    indexStart  uint64
    indexLength uint64
}

type SSTable struct {
    f             *os.File
    filePath      string
    tableMetaInfo MetaInfo // meta data
}
