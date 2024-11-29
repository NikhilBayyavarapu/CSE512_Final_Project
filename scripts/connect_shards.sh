#!/bin/bash

# Add shards to the cluster via Router-1
docker exec -it router-1 mongo --eval "sh.addShard('shard-1-replica-set/shard-1-node-a:27017');"
docker exec -it router-1 mongo --eval "sh.addShard('shard-2-replica-set/shard-2-node-a:27017');"
docker exec -it router-1 mongo --eval "sh.addShard('shard-3-replica-set/shard-3-node-a:27017');"

# Initialize the bank database
docker exec -it router-1 mongo --eval "use bank;"

# Enable sharding for the 'bank' database
docker exec -it router-1 mongo --eval "sh.enableSharding('bank');"

# Create collections
docker exec -it router-1 mongo --eval "use bank; db.createCollection('users');"
docker exec -it router-1 mongo --eval "use bank; db.createCollection('transactions');"

# Enable balancer and start it
docker exec -it router-1 mongo --eval "sh.setBalancerState(true);"
docker exec -it router-1 mongo --eval "sh.startBalancer();"

# Shard the 'users' collection on the hashed 'user_id' field
docker exec -it router-1 mongo --eval "sh.shardCollection('bank.users', { user_id: 'hashed' });"

# Shard the 'transactions' collection on the hashed 'transaction_id' field
docker exec -it router-1 mongo --eval "sh.shardCollection('bank.transactions', { transaction_id: 'hashed' });"