#!/usr/bin/env bash

env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o violet-etcd .

# upx --brute violet-etcd
