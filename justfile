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

up-test-server:
  docker compose -f ./test-server/docker-compose.yml up -d

down-test-server:
  docker compose -f ./test-server/docker-compose.yml down

launch port="8888":
  curl -X POST http://localhost:8000/servers/{{port}}/launch

down port="all":
  curl -X POST http://localhost:8000/servers/{{port}}/down
