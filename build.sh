#!/bin/bash

GOOS=linux GOARCH=amd64 go build -v -ldflags='-s -w' -o ./sreader ./cmd