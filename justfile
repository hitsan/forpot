build:
  go build -o ./bin/forpot ./cmd/forpot/main.go

run:
  go run ./cmd/forpot

clean:
  rm -fr ./bin/*

test module:
  go test {{module}}

test-all:
  go test ./...
