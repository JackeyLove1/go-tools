package main

import "fmt"

type User struct {
    Name string
    Age  int
    Tags map[string]string
}

func (u *User) String() string {
    return fmt.Sprintf("Name:%s, Age:%d,Tags:%v", u.Name, u.Age, u.Tags)
}

type UserOptions interface {
    Apply(user *User) error
}

type UserOption func(user *User) error

func NewUserWithOptions(opts ...UserOption) *User {
    user := &User{}
    for _, opt := range opts {
        err := opt(user)
        if err != nil {
            panic(err)
        }
    }
    return user
}

func WithName(name string) UserOption {
    return func(user *User) error {
        user.Name = name
        return nil
    }
}

func WithAge(age int) UserOption {
    return func(user *User) error {
        user.Age = age
        return nil
    }
}

func WithTag(key, value string) UserOption {
    return func(user *User) error {
        if user.Tags == nil {
            user.Tags = make(map[string]string)
        }
        user.Tags[key] = value
        return nil
    }
}

func main() {
    user := NewUserWithOptions(WithName("Jacky"), WithAge(18), WithTag("key1", "value1"))
    println(user.String())
}
