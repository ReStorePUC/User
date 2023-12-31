#!/bin/bash

cont=$(docker ps | grep mysql | awk '{print $1}')

docker cp "$(pwd)"/migrations "$cont":/

docker exec -it "$cont" bash -c "mysql < /migrations/$1.sql"