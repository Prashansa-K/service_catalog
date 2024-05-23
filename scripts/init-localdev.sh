#!/bin/bash

# export the default env VARS
export DB_PASSWORD="mysecretpassword"
export POSTGRES_DB="servicecatalog"
export API_AUTH_KEY="valid-key"
export DB_NAME="servicecatalog"
export JAEGER_SERVICE_NAME="servicecatalog"
export JAEGER_AGENT_HOST="localhost"
export JAEGER_AGENT_PORT="6831"


# check if a container name `postgres-db` is running
# if not, start it
docker ps -a | grep postgres-db > /dev/null
if [ $? -eq 0 ]; then
    if [ "$( docker container inspect -f '{{.State.Running}}' postgres-db )" == "true" ]; then
        echo "Container Postgres is already running"
    elif [ "$(docker ps -aq -f status=exited -f name=postgres-db)" ]; then
        echo "Container Postgres is stopped, starting it"
        docker start postgres-db
    fi
else
    # container doesn't exist, create it
    docker run -p 5432:5432  --name postgres-db -e POSTGRES_PASSWORD=${DB_PASSWORD} -e POSTGRES_DB=${POSTGRES_DB}  -d postgres
    # wait for the container to be ready
    sleep 5
    # apply the init migrations
    docker exec -i postgres-db /bin/bash -c "PGPASSWORD=${DB_PASSWORD} psql --username postgres ${POSTGRES_DB}" < $(PWD)/migrations/0.sql
    echo $(PWD)/migrations/0.sql
fi

# check the same for jaeger
docker ps -a | grep jaeger > /dev/null
if [ $? -eq 0 ]; then
    if [ "$( docker container inspect -f '{{.State.Running}}' jaeger )" == "true" ]; then
        echo "Container Jaeger is already running"
    elif [ "$(docker ps -aq -f status=exited -f name=jaeger)" ]; then
        echo "Container Jaeger is stopped, starting it"
        docker start jaeger
    fi
else
    # container doesn't exist, create it
    docker run -d --name jaeger -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p 5775:5775/udp -p ${JAEGER_AGENT_PORT}:${JAEGER_AGENT_PORT}/udp -p 6832:6832/udp -p 5778:5778 -p 16686:16686 -p 14268:14268 -p 9411:9411 jaegertracing/all-in-one:1.6
fi
