package LSM

import (
    "go-tools/db/LSM/ssTable"
)

type DataBase struct {
    TableTree *sstable.TableTree
}
