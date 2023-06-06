#!/bin/bash

rm -rf ./bin/linux
mkdir -p ./bin/linux

set -o xtrace

GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o ./bin/linux/jsonTester_386 ./cmd/jsonTester/main.go
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/linux/jsonTester_amd64 ./cmd/jsonTester/main.go
GOOS=linux GOARCH=arm go build -ldflags "-s -w" -o ./bin/linux/jsonTester_arm ./cmd/jsonTester/main.go
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o ./bin/linux/jsonTester_arm64 ./cmd/jsonTester/main.go
GOOS=linux GOARCH=loong64 go build -ldflags "-s -w" -o ./bin/linux/jsonTester_loong64 ./cmd/jsonTester/main.go
GOOS=linux GOARCH=mips go build -ldflags "-s -w" -o ./bin/linux/jsonTester_mips ./cmd/jsonTester/main.go
GOOS=linux GOARCH=mipsle go build -ldflags "-s -w" -o ./bin/linux/jsonTester_mipsle ./cmd/jsonTester/main.go
GOOS=linux GOARCH=mips64 go build -ldflags "-s -w" -o ./bin/linux/jsonTester_mips64 ./cmd/jsonTester/main.go
GOOS=linux GOARCH=mips64le go build -ldflags "-s -w" -o ./bin/linux/jsonTester_mips64le ./cmd/jsonTester/main.go
GOOS=linux GOARCH=ppc64 go build -ldflags "-s -w" -o ./bin/linux/jsonTester_ppc64 ./cmd/jsonTester/main.go
GOOS=linux GOARCH=ppc64le go build -ldflags "-s -w" -o ./bin/linux/jsonTester_ppc64le ./cmd/jsonTester/main.go
GOOS=linux GOARCH=riscv64 go build -ldflags "-s -w" -o ./bin/linux/jsonTester_riscv64 ./cmd/jsonTester/main.go
GOOS=linux GOARCH=s390x go build -ldflags "-s -w" -o ./bin/linux/jsonTester_s390x ./cmd/jsonTester/main.go

set +o xtrace