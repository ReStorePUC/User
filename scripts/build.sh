#!/bin/bash

echo "Compiling the API"
docker run -it --rm -v "$(pwd)":/go -e GOPATH= golang:1.21 sh -c "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o user ./cmd/main.go"

rm ./docker/user
mv ./user ./docker/
cp ./docker-config.yaml ./docker/config.yaml

docker build -t restore/user:"$1" docker/

if [[ ! $(docker service ls | grep rs_user) = "" ]]; then
  docker service update rs_user --image restore/user:$1
else
  docker stack deploy -c docker-compose.yaml rs
fi