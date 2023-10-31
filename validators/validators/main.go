package main

import (
    "fmt"

    "github.com/go-playground/validator/v10"
)

type User struct {
    Name string `validate:"contains=tom"`
    Age  int    `validate:"min=1"`
}

func main() {
    validate := validator.New()
    err := validate.Struct(User{Name: "tom", Age: 1})
    if err != nil {
        fmt.Println(err)
    }
}
