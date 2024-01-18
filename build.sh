#!/bin/bash
go mod tidy

binary=zigo

echo "building $binary-linux-amd64"
go env -w GOARCH=amd64
go build -ldflags="-s -w"
tar -czf $binary-linux-amd64.tar.gz $binary

echo "building $binary-darwin-amd64"
go env -w GOOS=darwin
go build -ldflags="-s -w"
tar -czf $binary-darwin-amd64.tar.gz $binary

echo "building $binary-windows-amd64"
go env -w GOOS=windows
go build -ldflags="-s -w"
zip $binary-windows-amd64.zip $binary.exe

echo "building $binary-darwin-arm64"
go env -w GOARCH=arm64
go env -w GOOS=darwin
go build -ldflags="-s -w"
tar -czf $binary-linux-arm64.tar.gz $binary

echo "building $binary-linux-arm64"
go env -w GOOS=linux
go build -ldflags="-s -w"
tar -czf $binary-darwin-arm64.tar.gz $binary
