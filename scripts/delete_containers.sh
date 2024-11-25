#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

echo "Deleting all containers related to the MongoDB sharded cluster setup..."

# List of container names to remove
containers=(
  "config-svr-1"
  "config-svr-2"
  "config-svr-3"
  "shard-1-node-a"
  "shard-1-node-b"
  "shard-1-node-c"
  "shard-2-node-a"
  "shard-2-node-b"
  "shard-2-node-c"
  "shard-3-node-a"
  "shard-3-node-b"
  "shard-3-node-c"
  "router-1"
  "router-2"
  "router-3"
)

# Loop through the container names and remove them if they exist
for container in "${containers[@]}"; do
  if docker ps -a --format '{{.Names}}' | grep -q "^${container}$"; then
    echo "Removing container: $container"
    docker rm -f "$container"
  else
    echo "Container $container does not exist. Skipping..."
  fi
done

echo "All specified containers have been removed."