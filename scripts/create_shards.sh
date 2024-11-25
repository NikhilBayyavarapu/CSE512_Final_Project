#!/bin/bash

# Function to create a shard replica set
create_shard_replica_set() {
  local shard_id=$1

  # Start shard nodes
  docker run -d --net mongo-shard-cluster --name shard-${shard_id}-node-a -p $((27110 + shard_id * 10 + 1)):27017 mongo:4.4 mongod --port 27017 --shardsvr --replSet shard-${shard_id}-replica-set
  docker run -d --net mongo-shard-cluster --name shard-${shard_id}-node-b -p $((27110 + shard_id * 10 + 2)):27017 mongo:4.4 mongod --port 27017 --shardsvr --replSet shard-${shard_id}-replica-set
  docker run -d --net mongo-shard-cluster --name shard-${shard_id}-node-c -p $((27110 + shard_id * 10 + 3)):27017 mongo:4.4 mongod --port 27017 --shardsvr --replSet shard-${shard_id}-replica-set

  # Wait for containers to start
  sleep 5

  # Initiate Shard Replica Set
  docker exec -it shard-${shard_id}-node-a mongo --eval "
rs.initiate({
  _id: 'shard-${shard_id}-replica-set',
  members: [
    { _id: 0, host: 'shard-${shard_id}-node-a:27017' },
    { _id: 1, host: 'shard-${shard_id}-node-b:27017' },
    { _id: 2, host: 'shard-${shard_id}-node-c:27017' }
  ]
});
"
}

# Create three shard replica sets
create_shard_replica_set 1
create_shard_replica_set 2
create_shard_replica_set 3