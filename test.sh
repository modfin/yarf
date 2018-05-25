#!/bin/bash


echo "### Unit tests ###"
go test -v -race -count 1 . | awk '{print "\t" $0}'


echo "### Integration tests ###"
echo "  --- Starting NATS through docker"
docker run -d --name nats-test -p 4222:4222 -p 6222:6222 -p 8222:8222 nats > /dev/null

echo "  --- Starting tests"
go test -v -race -count 1 ./_integration_test/ | awk '{print "\t" $0}'

echo "  --- Killing NATS"
docker kill nats-test > /dev/null
docker rm nats-test > /dev/null
