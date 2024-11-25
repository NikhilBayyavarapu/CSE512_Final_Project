#!/bin/bash

# Start Router Instances
docker run -d --net mongo-shard-cluster --name router-1 -p 27151:27017 mongo:4.4 mongos \
--port 27017 \
--configdb config-svr-replica-set/config-svr-1:27017,config-svr-2:27017,config-svr-3:27017 \
--bind_ip_all

docker run -d --net mongo-shard-cluster --name router-2 -p 27152:27017 mongo:4.4 mongos \
--port 27017 \
--configdb config-svr-replica-set/config-svr-1:27017,config-svr-2:27017,config-svr-3:27017 \
--bind_ip_all

docker run -d --net mongo-shard-cluster --name router-3 -p 27153:27017 mongo:4.4 mongos \
--port 27017 \
--configdb config-svr-replica-set/config-svr-1:27017,config-svr-2:27017,config-svr-3:27017 \
--bind_ip_all

# Wait for routers to start
sleep 5