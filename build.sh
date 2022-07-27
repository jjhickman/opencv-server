#!/bin/bash
now=$(date +'%Y-%m-%d_%T')
go build -o ./build/telescope -ldflags="-s -w -X main.buildTime=$now -X main.version=$1" -tags static ./cmd/telescope