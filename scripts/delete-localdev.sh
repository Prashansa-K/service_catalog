#!/bin/bash

docker rm -f postgres-db || true
docker rm -f jaeger || true
