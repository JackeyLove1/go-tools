package main

import "time"

type OP int

const (
    INSET OP = iota
    UPDATE
    DELETE
)

func (op OP) String() string {
    return [...]string{
        "INSERT",
        "UPDATE",
        "DELETE",
    }[op]
}

func main() {
    timeStamp := time.Now().Format("2006-01-02 15:04:05")
    println(timeStamp)
    println(INSET.String())
}
