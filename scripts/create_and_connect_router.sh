#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Check if router ID is provided as an argument
if [ -z "$1" ]; then
  echo "Usage: $0 <router-id>"
  exit 1
fi

# Router ID from input
ROUTER_ID=$1

# Port for the new router (calculated dynamically based on router ID)
ROUTER_PORT=$((27150 + ROUTER_ID))

echo "Creating Router: router-${ROUTER_ID} on port ${ROUTER_PORT}"

# Start the new router (mongos)
docker run -d --net mongo-shard-cluster --name router-${ROUTER_ID} -p ${ROUTER_PORT}:27017 mongo:4.4 mongos \
  --port 27017 \
  --configdb config-svr-replica-set/config-svr-1:27017,config-svr-2:27017,config-svr-3:27017 \
  --bind_ip_all

echo "Router router-${ROUTER_ID} has been successfully created and connected to the config servers!"