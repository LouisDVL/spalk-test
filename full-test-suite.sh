#!/bin/bash
go mod tidy
go build index.go
./index < ./tests/test_success.ts
./index < ./tests/test_failure.ts