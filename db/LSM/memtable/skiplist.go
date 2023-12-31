// Copyright 2011 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

// Package skiplist implement skip list data structure.
// See wikipedia for more details about this data structure. http://en.wikipedia.org/wiki/Skip_list
//
// Skip list is basically an ordered map.
//
// Here is a sample to use this package.
//
//     // Creates a new skip list and restricts key type to int-like types.
//     list := skiplist.New(skiplist.Int)
//
//     // Adds some values for keys.
//     list.Set(20, "Hello")
//     list.Set(10, "World")
//     list.Set(40, true) // Value type is not restricted.
//     list.Set(40, 1000) // Replace the of an existing Element[K, V].
//
//     // Finds Element[K, V]s.
//     e := list.Get(10)           // Returns the Element[K, V] with the key.
//     _ = e.Value.(string)
//     v, ok := list.GetValue(20)  // Directly get value of the Element[K, V]. If the key is not found, ok is false.
//     v2 := list.MustGetValue(10) // Directly get value of the Element[K, V]. Panic if the key is not found.
//     notFound := list.Get(15)    // Returns nil if the key is not found.
//
//     // Removes an Element[K, V] and gets removed Element[K, V].
//     old := list.Remove(40)
//     notFound := list.Remove(-20) // Returns nil if the key is not found.
//
//     // Initializes the list again to clean up all Element[K, V]s in the list.
//     list.Init()
package memtable

import (
    "fmt"
    "math/rand"
    "time"
)

// DefaultMaxLevel is the default level for all newly created skip lists.
// It can be changed globally. Changing it will not affect existing lists.
// And all skip lists can update max level after creation through `SetMaxLevel()` method.
var DefaultMaxLevel = 48

// preallocDefaultMaxLevel is a constant to alloc memory on stack when Set new Element[K, V].
const preallocDefaultMaxLevel = 48

// SkipList is the header of a skip list.
type SkipList[K comparable, V any] struct {
    elementHeader[K, V]

    comparable Comparable
    rand       *rand.Rand

    maxLevel int
    length   int
    back     *Element[K, V]
}

// New creates a new skip list with comparable to compare keys.
//
// There are lots of pre-defined strict-typed keys like Int, Float64, String, etc.
// We can create custom comparable by implementing Comparable interface.
func New[K comparable, V any](comparable Comparable) *SkipList[K, V] {
    if DefaultMaxLevel <= 0 {
        panic("skiplist default level must not be zero or negative")
    }

    source := rand.NewSource(time.Now().UnixNano())
    return &SkipList[K, V]{
        elementHeader: elementHeader[K, V]{
            levels: make([]*Element[K, V], DefaultMaxLevel),
        },

        comparable: comparable,
        rand:       rand.New(source),

        maxLevel: DefaultMaxLevel,
    }
}

// Init resets the list and discards all existing Element[K, V]s.
func (list *SkipList[K, V]) Init() *SkipList[K, V] {
    list.back = nil
    list.length = 0
    list.levels = make([]*Element[K, V], len(list.levels))
    return list
}

// SetRandSource sets a new rand source.
//
// Skiplist uses global rand defined in math/rand by default.
// The default rand acquires a global mutex before generating any number.
// It's not necessary if the skiplist is well protected by caller.
func (list *SkipList[K, V]) SetRandSource(source rand.Source) {
    list.rand = rand.New(source)
}

// Front returns the first Element[K, V].
//
// The complexity is O(1).
func (list *SkipList[K, V]) Front() (front *Element[K, V]) {
    return list.levels[0]
}

// Back returns the last Element[K, V].
//
// The complexity is O(1).
func (list *SkipList[K, V]) Back() *Element[K, V] {
    return list.back
}

// Len returns Element[K, V] count in this list.
//
// The complexity is O(1).
func (list *SkipList[K, V]) Len() int {
    return list.length
}

// Set sets value for the key.
// If the key exists, updates Element[K, V]'s value.
// Returns the Element[K, V] holding the key and value.
//
// The complexity is O(log(N)).
func (list *SkipList[K, V]) Set(key K, value V) (elem *Element[K, V]) {
    score := list.calcScore(key)

    // Happy path for empty list.
    if list.length == 0 {
        level := list.randLevel()
        elem = newElement[K, V](list, level, score, key, value)

        for i := 0; i < level; i++ {
            list.levels[i] = elem
        }

        list.back = elem
        list.length++
        return
    }

    // Find out previous Element[K, V]s on every possible levels.
    max := len(list.levels)
    prevHeader := &list.elementHeader

    var maxStaticAllocElemHeaders [preallocDefaultMaxLevel]*elementHeader[K, V]
    var prevElemHeaders []*elementHeader[K, V]

    if max <= preallocDefaultMaxLevel {
        prevElemHeaders = maxStaticAllocElemHeaders[:max]
    } else {
        prevElemHeaders = make([]*elementHeader[K, V], max)
    }

    for i := max - 1; i >= 0; {
        prevElemHeaders[i] = prevHeader

        for next := prevHeader.levels[i]; next != nil; next = prevHeader.levels[i] {
            if comp := list.compare(score, key, next); comp <= 0 {
                // Find the elem with the same key.
                // Update value and return the elem.
                if comp == 0 {
                    elem = next
                    elem.Value = value
                    return
                }

                break
            }

            prevHeader = &next.elementHeader
            prevElemHeaders[i] = prevHeader
        }

        // Skip levels if they point to the same Element[K, V] as topLevel.
        topLevel := prevHeader.levels[i]

        for i--; i >= 0 && prevHeader.levels[i] == topLevel; i-- {
            prevElemHeaders[i] = prevHeader
        }
    }

    // Create a new Element[K, V].
    level := list.randLevel()
    elem = newElement[K, V](list, level, score, key, value)

    // Set up prev Element[K, V].
    if prev := prevElemHeaders[0]; prev != &list.elementHeader {
        elem.prev = prev.Element()
    }

    // Set up prevTopLevel.
    if prev := prevElemHeaders[level-1]; prev != &list.elementHeader {
        elem.prevTopLevel = prev.Element()
    }

    // Set up levels.
    for i := 0; i < level; i++ {
        elem.levels[i] = prevElemHeaders[i].levels[i]
        prevElemHeaders[i].levels[i] = elem
    }

    // Find out the largest level with next Element[K, V].
    largestLevel := 0

    for i := level - 1; i >= 0; i-- {
        if elem.levels[i] != nil {
            largestLevel = i + 1
            break
        }
    }

    // Adjust prev and prevTopLevel of next Element[K, V]s.
    if next := elem.levels[0]; next != nil {
        next.prev = elem
    }

    for i := 0; i < largestLevel; {
        next := elem.levels[i]
        nextLevel := next.Level()

        if nextLevel <= level {
            next.prevTopLevel = elem
        }

        i = nextLevel
    }

    // If the elem is the last Element[K, V], set it as back.
    if elem.Next() == nil {
        list.back = elem
    }

    list.length++
    return
}

func (list *SkipList[K, V]) findNext(start *Element[K, V], score float64, key K) (elem *Element[K, V]) {
    if list.length == 0 {
        return
    }

    if start == nil && list.compare(score, key, list.Front()) <= 0 {
        elem = list.Front()
        return
    }
    if start != nil && list.compare(score, key, start) <= 0 {
        elem = start
        return
    }
    if list.compare(score, key, list.Back()) > 0 {
        return
    }

    var prevHeader *elementHeader[K, V]
    if start == nil {
        prevHeader = &list.elementHeader
    } else {
        prevHeader = &start.elementHeader
    }
    i := len(prevHeader.levels) - 1

    // Find out previous Element[K, V]s on every possible levels.
    for i >= 0 {
        for next := prevHeader.levels[i]; next != nil; next = prevHeader.levels[i] {
            if comp := list.compare(score, key, next); comp <= 0 {
                elem = next
                if comp == 0 {
                    return
                }

                break
            }

            prevHeader = &next.elementHeader
        }

        topLevel := prevHeader.levels[i]

        // Skip levels if they point to the same Element[K, V] as topLevel.
        for i--; i >= 0 && prevHeader.levels[i] == topLevel; i-- {
        }
    }

    return
}

// FindNext returns the first Element[K, V] after start that is greater or equal to key.
// If start is greater or equal to key, returns start.
// If there is no such Element[K, V], returns nil.
// If start is nil, find Element[K, V] from front.
//
// The complexity is O(log(N)).
func (list *SkipList[K, V]) FindNext(start *Element[K, V], key K) (elem *Element[K, V]) {
    return list.findNext(start, list.calcScore(key), key)
}

// Find returns the first Element[K, V] that is greater or equal to key.
// It's short hand for FindNext(nil, key).
//
// The complexity is O(log(N)).
func (list *SkipList[K, V]) Find(key K) (elem *Element[K, V]) {
    return list.FindNext(nil, key)
}

// Get returns an Element[K, V] with the key.
// If the key is not found, returns nil.
//
// The complexity is O(log(N)).
func (list *SkipList[K, V]) Get(key K) (elem *Element[K, V]) {
    score := list.calcScore(key)

    firstElem := list.findNext(nil, score, key)
    if firstElem == nil {
        return
    }

    if list.compare(score, key, firstElem) != 0 {
        return
    }

    elem = firstElem
    return
}

// GetValue returns value of the Element[K, V] with the key.
// It's short hand for Get().Value.
//
// The complexity is O(log(N)).
func (list *SkipList[K, V]) GetValue(key K) (val V, ok bool) {
    element := list.Get(key)

    if element == nil {
        return
    }

    val = element.Value
    ok = true
    return
}

// MustGetValue returns value of the Element[K, V] with the key.
// It will panic if the key doesn't exist in the list.
//
// The complexity is O(log(N)).
func (list *SkipList[K, V]) MustGetValue(key K) interface{} {
    element := list.Get(key)

    if element == nil {
        panic(fmt.Errorf("skiplist: cannot find key `%v` in skiplist", key))
    }

    return element.Value
}

// Remove removes an Element[K, V].
// Returns removed Element[K, V] pointer if found, nil if it's not found.
//
// The complexity is O(log(N)).
func (list *SkipList[K, V]) Remove(key K) (elem *Element[K, V]) {
    elem = list.Get(key)

    if elem == nil {
        return
    }

    list.RemoveElement(elem)
    return
}

// RemoveFront removes front Element[K, V] node and returns the removed Element[K, V].
//
// The complexity is O(1).
func (list *SkipList[K, V]) RemoveFront() (front *Element[K, V]) {
    if list.length == 0 {
        return
    }

    front = list.Front()
    list.RemoveElement(front)
    return
}

// RemoveBack removes back Element[K, V] node and returns the removed Element[K, V].
//
// The complexity is O(log(N)).
func (list *SkipList[K, V]) RemoveBack() (back *Element[K, V]) {
    if list.length == 0 {
        return
    }

    back = list.back
    list.RemoveElement(back)
    return
}

// RemoveElement[K, V] removes the elem from the list.
//
// The complexity is O(log(N)).
func (list *SkipList[K, V]) RemoveElement(elem *Element[K, V]) {
    if elem == nil || elem.list != list {
        return
    }

    level := elem.Level()

    // Find out all previous Element[K, V]s.
    max := 0
    prevElems := make([]*Element[K, V], level)
    prev := elem.prev

    for prev != nil && max < level {
        prevLevel := len(prev.levels)

        for ; max < prevLevel && max < level; max++ {
            prevElems[max] = prev
        }

        for prev = prev.prevTopLevel; prev != nil && prev.Level() == prevLevel; prev = prev.prevTopLevel {
        }
    }

    // Adjust prev Element[K, V]s which point to elem directly.
    for i := 0; i < max; i++ {
        prevElems[i].levels[i] = elem.levels[i]
    }

    for i := max; i < level; i++ {
        list.levels[i] = elem.levels[i]
    }

    // Adjust prev and prevTopLevel of next Element[K, V]s.
    if next := elem.Next(); next != nil {
        next.prev = elem.prev
    }

    for i := 0; i < level; {
        next := elem.levels[i]

        if next == nil || next.prevTopLevel != elem {
            break
        }

        i = next.Level()
        next.prevTopLevel = prevElems[i-1]
    }

    // Adjust list.Back() if necessary.
    if list.back == elem {
        list.back = elem.prev
    }

    list.length--
    elem.reset()
}

// MaxLevel returns current max level value.
func (list *SkipList[K, V]) MaxLevel() int {
    return list.maxLevel
}

// SetMaxLevel changes skip list max level.
// If level is not greater than 0, just panic.
func (list *SkipList[K, V]) SetMaxLevel(level int) (old int) {
    if level <= 0 {
        panic(fmt.Errorf("skiplist: level must be larger than 0 (current is %v)", level))
    }

    list.maxLevel = level
    old = len(list.levels)

    if level == old {
        return
    }

    if old > level {
        for i := old - 1; i >= level; i-- {
            if list.levels[i] != nil {
                level = i
                break
            }
        }

        list.levels = list.levels[:level]
        return
    }

    if level <= cap(list.levels) {
        list.levels = list.levels[:level]
        return
    }

    levels := make([]*Element[K, V], level)
    copy(levels, list.levels)
    list.levels = levels
    return
}

func (list *SkipList[K, V]) randLevel() int {
    estimated := list.maxLevel
    const prob = 1 << 30 // Half of 2^31.
    rand := list.rand
    i := 1

    for ; i < estimated; i++ {
        if rand.Int31() < prob {
            break
        }
    }

    return i
}

// compare compares value of two Element[K, V]s and returns -1, 0 and 1.
func (list *SkipList[K, V]) compare(score float64, key interface{}, rhs *Element[K, V]) int {
    if score != rhs.score {
        if score > rhs.score {
            return 1
        } else if score < rhs.score {
            return -1
        }

        return 0
    }

    return list.comparable.Compare(key, rhs.key)
}

func (list *SkipList[K, V]) calcScore(key any) (score float64) {
    score = list.comparable.CalcScore(key)
    return
}
