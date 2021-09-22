#!/bin/bash

GOOS=linux GOARCH=arm64 go build -v -ldflags='-s -w' -o ./sreader ./cmd