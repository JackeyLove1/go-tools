#!/bin/bash
cur=$(pwd)
mkdir -p "$cur"/idl/student_service
protoc --go_out=./idl/student_service  --go_opt=paths=source_relative --go-grpc_out=./idl/student_service --go-grpc_opt=paths=source_relative student_service.proto