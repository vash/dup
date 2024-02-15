dup:
  go mod download && GOOS=linux GOARCH=amd64 go build -a -ldflags="-s -w" -o ./bin/kubectl-dup .
clean:
  rm -rf bin
test:
  go test ...
