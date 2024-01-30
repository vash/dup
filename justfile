dup:
  go mod download && GOOS=linux GOARCH=amd64 go build -a -ldflags="-s -w" -o ./bin/kubectl-dup .
darwin-intel:
  go mod download && GOOS=darwin GOARCH=amd64 go build -a -ldflags="-s -w" -o ./bin/kubectl-dup .
darwin-apple:
  go mod download && GOOS=darwin GOARCH=arm64 go build -a -ldflags="-s -w" -o ./bin/kubectl-dup .
windows:
  go mod download && GOOS=windows GOARCH=amd64 go build -a -ldflags="-s -w" -o ./bin/kubectl-dup .
package:
  mkdir -p dist
  rm -rf dist/*

  go mod download && GOOS=linux GOARCH=amd64 go build -a -ldflags="-s -w" -o ./bin/kubectl-dup .
  tar -cf dist/linux_amd64.tar.gz ./bin/kubectl-dup

  go mod download && GOOS=darwin GOARCH=amd64 go build -a -ldflags="-s -w" -o ./bin/kubectl-dup .
  tar -cf dist/darwin_amd64.tar.gz ./bin/kubectl-dup

  go mod download && GOOS=darwin GOARCH=arm64 go build -a -ldflags="-s -w" -o ./bin/kubectl-dup .
  tar -cf dist/darwin_arm64.tar.gz ./bin/kubectl-dup

  go mod download && GOOS=windows GOARCH=amd64 go build -a -ldflags="-s -w" -o ./bin/kubectl-dup .
  tar -cf dist/windows_amd64.tar.gz ./bin/kubectl-dup

clean:
  rm -rf bin
test:
  go test ...
