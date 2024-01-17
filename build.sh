go mod tidy

echo "building zigo-linux-amd64"
go env -w GOARCH=amd64
go build -ldflags="-s -w" -o "dist/zigo-linux-amd64"

echo "building zigo-darwin-amd64"
go env -w GOOS=darwin
go build -ldflags="-s -w" -o "dist/zigo-darwin-amd64"

echo "building zigo-windows-amd64"
go env -w GOOS=windows
go build -ldflags="-s -w" -o "dist/zigo-windows-amd64.exe"

echo "building zigo-darwin-arm64"
go env -w GOARCH=arm64
go env -w GOOS=darwin
go build -ldflags="-s -w" -o "dist/zigo-darwin-arm64"

echo "building zigo-linux-arm64"
go env -w GOOS=linux
go build -ldflags="-s -w" -o "dist/zigo-linux-arm64"

cd dist
tar -czf zigo-linux-amd64.tar.gz zigo-linux-amd64
tar -czf zigo-darwin-amd64.tar.gz zigo-darwin-amd64
tar -czf zigo-linux-arm64.tar.gz zigo-linux-arm64
tar -czf zigo-darwin-arm64.tar.gz zigo-darwin-arm64
zip zigo-windows-amd64.zip zigo-windows-amd64.exe
