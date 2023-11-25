package main

import (
    "context"

    "go-tools/practice/v1/v2/idl/student_service"
    "google.golang.org/grpc"
)

func main() {
    conn, err := grpc.Dial("localhost:2346", grpc.WithInsecure())
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    client := student_service.NewStudentServiceClient(conn)
    resp, err := client.GetStudentInfo(context.TODO(), &student_service.Request{
        StudentId: "hello, world!",
    })
    if err != nil {
        panic(err)
    }
    println(resp.String())
}
