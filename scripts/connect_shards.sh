#!/bin/bash

# Add shards to the cluster via Router-1
docker exec -it router-1 mongo --eval '
sh.addShard("shard-1-replica-set/shard-1-node-a:27017");
sh.addShard("shard-2-replica-set/shard-2-node-a:27017");
sh.addShard("shard-3-replica-set/shard-3-node-a:27017");
'