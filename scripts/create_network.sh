#!/bin/bash

# Check if the Docker network already exists
if docker network inspect mongo-shard-cluster >/dev/null 2>&1; then
  echo "Docker network 'mongo-shard-cluster' already exists. Removing it..."
  docker network rm mongo-shard-cluster
fi

# Create a new Docker network
echo "Creating Docker network 'mongo-shard-cluster'..."
docker network create mongo-shard-cluster