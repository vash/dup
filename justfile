duplipod:
  go mod download && GOOS=linux GOARCH=amd64 go build -a -ldflags="-s -w" -o ./bin/duplipod ./cmd
local:
  go mod download && GOOS=linux GOARCH=amd64 go build -a -o ./bin/duplipod-local ./cmd
clean:
  rm -rf bin
test:
  go test ...
