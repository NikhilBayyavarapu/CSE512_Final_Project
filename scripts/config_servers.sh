#!/bin/bash

# Start Config Servers
docker run -d --net mongo-shard-cluster --name config-svr-1 -p 27101:27017 mongo:4.4 mongod --port 27017 --configsvr --replSet config-svr-replica-set
docker run -d --net mongo-shard-cluster --name config-svr-2 -p 27102:27017 mongo:4.4 mongod --port 27017 --configsvr --replSet config-svr-replica-set
docker run -d --net mongo-shard-cluster --name config-svr-3 -p 27103:27017 mongo:4.4 mongod --port 27017 --configsvr --replSet config-svr-replica-set

# Wait for containers to start
sleep 5

# Initiate Config Server Replica Set
docker exec -it config-svr-1 mongo --eval '
rs.initiate({
  _id: "config-svr-replica-set",
  configsvr: true,
  members: [
    { _id: 0, host: "config-svr-1:27017" },
    { _id: 1, host: "config-svr-2:27017" },
    { _id: 2, host: "config-svr-3:27017" }
  ]
});
'