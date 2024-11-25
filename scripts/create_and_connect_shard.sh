#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Define color codes for aesthetics
RESET='\033[0m'
BOLD='\033[1m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
RED='\033[0;31m'

# Check if shard ID is provided as an argument
if [ -z "$1" ]; then
  echo -e "${RED}${BOLD}Usage: $0 <shard-id>${RESET}"
  exit 1
fi

# Shard ID from input
SHARD_ID=$1

echo -e "\n${CYAN}${BOLD}Creating Shard Replica Set: shard-${SHARD_ID}-replica-set${RESET}"

# Start shard nodes
echo -e "${BLUE}Starting shard nodes...${RESET}"

docker run -d --net mongo-shard-cluster --name shard-${SHARD_ID}-node-a -p $((27110 + SHARD_ID * 10 + 1)):27017 mongo:4.4 mongod --port 27017 --shardsvr --replSet shard-${SHARD_ID}-replica-set
docker run -d --net mongo-shard-cluster --name shard-${SHARD_ID}-node-b -p $((27110 + SHARD_ID * 10 + 2)):27017 mongo:4.4 mongod --port 27017 --shardsvr --replSet shard-${SHARD_ID}-replica-set
docker run -d --net mongo-shard-cluster --name shard-${SHARD_ID}-node-c -p $((27110 + SHARD_ID * 10 + 3)):27017 mongo:4.4 mongod --port 27017 --shardsvr --replSet shard-${SHARD_ID}-replica-set

# Wait for containers to start
echo -e "\n${YELLOW}Waiting for shard nodes to start...${RESET}"
sleep 5

echo -e "\n${GREEN}${BOLD}Initiating Shard Replica Set: shard-${SHARD_ID}-replica-set${RESET}"

# Initiate the replica set
docker exec -it shard-${SHARD_ID}-node-a mongo --eval "
rs.initiate({
  _id: 'shard-${SHARD_ID}-replica-set',
  members: [
    { _id: 0, host: 'shard-${SHARD_ID}-node-a:27017' },
    { _id: 1, host: 'shard-${SHARD_ID}-node-b:27017' },
    { _id: 2, host: 'shard-${SHARD_ID}-node-c:27017' }
  ]
});
rs.status();
"

# Connect the shard to the cluster via Router-1
echo -e "\n${BLUE}Connecting Shard Replica Set: shard-${SHARD_ID}-replica-set to Routers${RESET}"
docker exec -it router-1 mongo --eval "
sh.addShard('shard-${SHARD_ID}-replica-set/shard-${SHARD_ID}-node-a:27017');
"

echo -e "\n${GREEN}${BOLD}Shard ${SHARD_ID} has been successfully created and connected to the cluster!${RESET}"

