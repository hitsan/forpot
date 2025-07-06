build:
  go build -o ./bin/forpot ./cmd/forpot/main.go

run *args:
  go run ./cmd/forpot {{args}}

clean:
  rm -fr ./bin/*

test module="all":
  if [ "{{module}}" == "all" ]; then \
    go test ./...; \
  else \
    go test ./{{module}}; \
  fi

test-server arg="":
  if [ "{{arg}}" == "up" ]; then \
    docker compose -f ./test-server/docker-compose.yml up -d; \
  elif [ "{{arg}}" == "down" ]; then \
    docker compose -f ./test-server/docker-compose.yml down; \
  else \
    echo "Illigal argument!"; \
  fi

launch port="8888":
  curl -X POST http://localhost:8000/servers/{{port}}/launch

get port="8888":
  curl -X GET http://localhost:{{port}}

down port="all":
  curl -X POST http://localhost:8000/servers/{{port}}/down
