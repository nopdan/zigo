echo "building zigo-windows-amd64"
go build -ldflags="-s -w" -o "build/zigo-windows-amd64/zigo.exe"

echo "building zigo-linux-amd64"
go env -w GOOS=linux
go build -ldflags="-s -w" -o "build/zigo-linux-amd64/zigo"

echo "building zigo-darwin-amd64"
go env -w GOOS=darwin
go build -ldflags="-s -w" -o "build/zigo-darwin-amd64/zigo"

echo "building zigo-windows-arm64"
go env -w GOARCH=arm64
go build -ldflags="-s -w" -o "build/zigo-darwin-arm64/zigo"

echo "building zigo-linux-arm64"
go env -w GOOS=linux
go build -ldflags="-s -w" -o "build/zigo-linux-arm64/zigo"

# reset go env
go env -w GOARCH=amd64
go env -w GOOS=windows
