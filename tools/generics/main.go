package main

func MapKeys[K comparable, V any](m map[K]V) []K {
    r := make([]K, 0, len(m))
    for k, _ := range m {
        r = append(r, k)
    }
    return r
}

type Element[T any] struct {
    next *Element[T]
    val  T
}

type List[T any] struct {
    head, tail *Element[T]
}

func (lst *List[T]) Push(v T) {

}

func main() {

}
