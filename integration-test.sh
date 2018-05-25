#!/bin/bash



docker run -d --name nats-test -p 4222:4222 -p 6222:6222 -p 8222:8222 nats

go test -v -race -count 1 ./_integration_test/

docker kill nats-test
docker rm nats-test
