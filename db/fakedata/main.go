package main

import "github.com/brianvoe/gofakeit/v6"

func main() {
    println(gofakeit.Name())
    println(gofakeit.Email())
    println(gofakeit.Date().String())
    println(gofakeit.Username())
    println(gofakeit.Second())
}
