package main

import (
    "context"
    "math/rand"
    "net"

    "github.com/brianvoe/gofakeit/v6"
    student_service "go-tools/practice/v1/v2/idl/student_service"
    "google.golang.org/grpc"
)

var fake *gofakeit.Faker

func init() {
    fake = gofakeit.New(rand.Int63())
}

type StudentService struct {
    student_service.UnimplementedStudentServiceServer
}

func (s *StudentService) GetStudentInfo(ctx context.Context, req *student_service.Request) (*student_service.Student, error) {
    defer func() {
        if err := recover(); err != nil {
            println("panic in run get student info, err:", err)
        }
    }()
    studentId := req.GetStudentId()
    return &student_service.Student{
        Id:        studentId,
        Name:      fake.Name(),
        Age:       fake.Int32(),
        Height:    fake.Float32(),
        Locations: []string{"Beijing", "Shanghai", "Chengdu", "Guangzhou", "Shenzhen"},
        Scores:    make(map[string]float32),
    }, nil
}

func main() {
    list, err := net.Listen("tcp", ":2346")
    if err != nil {
        panic(err)
    }
    server := grpc.NewServer()
    student_service.RegisterStudentServiceServer(server, &StudentService{})
    err = server.Serve(list)
    if err != nil {
        panic(err)
    }
}
