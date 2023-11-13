package main

import (
    "fmt"

    "github.com/huandu/skiplist"
)

func main() {
    // Create a skip list with int key.
    list := skiplist.New(skiplist.Int)

    // Add some values. Value can be anything.
    list.Set(12, "hello world")
    list.Set(34, 56)
    list.Set(78, 90.12)

    // Get element by index.
    elem := list.Get(34)                // Value is stored in elem.Value.
    fmt.Println(elem.Value)             // Output: 56
    next := elem.Next()                 // Get next element.
    prev := next.Prev()                 // Get previous element.
    fmt.Println(next.Value, prev.Value) // Output: 90.12    56

    // Or, directly get value just like a map
    val, ok := list.GetValue(34)
    fmt.Println(val, ok) // Output: 56  true

    // Find first elements with score greater or equal to key
    foundElem := list.Find(30)
    fmt.Println(foundElem.Key(), foundElem.Value) // Output: 34 56

    // Remove an element for key.
    // list.Remove(34)
    println("iterator")
    for iter := list.Front(); iter != nil; iter = iter.Next() {
        fmt.Println(iter.Key(), iter.Value)
    }
}
