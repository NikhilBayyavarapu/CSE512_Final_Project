## Create a Docker Network
```
docker network create mongo-shard-cluster
```

## Config Servers

Config servers store metadata of mongo shards. This metadata will be useful for routers when deciding on which shard to direct the query to.

```
docker run -d --net mongo-shard-cluster --name config-svr-1 -p 27101:27017 mongo:4.4 mongod --port 27017 --configsvr --replSet config-svr-replica-set
```

```
docker run -d --net mongo-shard-cluster --name config-svr-2 -p 27102:27017 mongo:4.4 mongod --port 27017 --configsvr --replSet config-svr-replica-set
```

```
docker run -d --net mongo-shard-cluster --name config-svr-3 -p 27103:27017 mongo:4.4 mongod --port 27017 --configsvr --replSet config-svr-replica-set
```

We have mentioned that all these 3 config servers belong to a replicaSet named *config-svr-replica-set*.

So lets proceed and initiate this replicaSet.

First access mongo shell of any config server:
```
docker exec -it config-svr-1 mongo
```

Then initiate the replicaSet with the following command:
```
rs.initiate({
    _id: "config-svr-replica-set",
    configsvr: true,
    members: [
        { _id: 0, host: "config-svr-1:27017" },
        { _id: 1, host: "config-svr-2:27017" },
        { _id: 2, host: "config-svr-3:27017" }
    ]
})
```
Since the above command is run inside docker, all of the servers are running on port 27017. But outside of docker, these ports are exposed on 27101,27102,27103 respectively.

In host, specify the address of the previously created config servers.

To check the status of the replicaSet. Run the following command:

```
rs.status()
```

One of the 3 config servers will be PRIMARY and the others will be SECONDARY. Mongo will ping consistenly to ensure PRIMARY server will always be available. 

If config-svr-1 does not respond, mongo will automatically promote one of the other two servers to PRIMARY status.

## Creating Shards

Creating shards follows a similar approach to that of creating config servers.

Our goal here is to create 3 shard replicaSets. And in each replicaSet, we will have 3 shards. One of the 3 shards will be in each replicaSet will be PRIMARY and the other two will be SECONDARY. Data will be replicated across all the shards within a replicaSet.

When inserting data, our router will decide in shard to insert the data. We can specify a field of our data to router. Then router hashes the value in this field and inserts data in appropriate shard.

### Shard Replica Set 1
```
docker run -d --net mongo-shard-cluster --name shard-1-node-a -p 27111:27017 mongo:4.4 mongod --port 27017 --shardsvr --replSet shard-1-replica-set
```

```
docker run -d --net mongo-shard-cluster --name shard-1-node-b -p 27112:27017 mongo:4.4 mongod --port 27017 --shardsvr --replSet shard-1-replica-set
```

```
docker run -d --net mongo-shard-cluster --name shard-1-node-c -p 27113:27017 mongo:4.4 mongod --port 27017 --shardsvr --replSet shard-1-replica-set
```

Similarly create 2 more shard replicaSets each with 3 nodes inside them.

Now we have to initiate *shard-1-replica-set*.

First access mongo shell of shard-1-node(a/b/c):

```
docker exec -it shard-1-node-a mongo
```

Then initiate the replicaSet:
```
rs.initiate({
    _id: "shard-3-replica-set",
    members: [
        { _id: 0, host: "shard-3-node-a:27017" },
        { _id: 1, host: "shard-3-node-b:27017" },
        { _id: 2, host: "shard-3-node-c:27017" }
    ]
})
```

Check status using:
```
rs.status()
```

Do not forget to create 2 more shard-replicaSets with 3 shards each.


## Creating Routers

Routers will be the entrypoint for our client applications.

### Router 1
```
docker run -d --net mongo-shard-cluster --name router-1 -p 27141:27017 mongo:4.4 mongos --port 27017 --configdb config-svr-replica-set/config-svr-1:27017,config-svr-2:27017,config-svr-3:27017 --bind_ip_all
```

### Router 2
```
docker run -d --net mongo-shard-cluster --name router-2 -p 27142:27017 mongo:4.4 mongos --port 27017 --configdb config-svr-replica-set/config-svr-1:27017,config-svr-2:27017,config-svr-3:27017 --bind_ip_all
```

Here we also bind our config servers to router. This way router will have access to metadata of shards and can redirect query to appropriate shard.

Now we have to connected config servers and routers. Let us proceed and connect our shards to this architecture.

## Connect Shards to Routers and Config Servers.

First access mongo shell of any router:
```
docker exec -it router-1 mongo
```

Add shards:

```
sh.addShard("shard-1-replica-set/shard-1-node-a:27017", "shard-1-replica-set/shard-1-node-b:27017", "shard-1-replica-set/shard-1-node-c:27017")
```

```
sh.addShard("shard-2-replica-set/shard-2-node-a:27017", "shard-2-replica-set/shard-2-node-b:27017", "shard-2-replica-set/shard-2-node-c:27017")
```

```
sh.addShard("shard-3-replica-set/shard-3-node-a:27017", "shard-3-replica-set/shard-3-node-b:27017", "shard-3-replica-set/shard-3-node-c:27017")
```

Why does running this command on a single router does the job ?

*When creating routers, we linked them to our config servers. So, now both of our routers to config servers. When we run sh.addShard() command , it contacts with config servers and stores data there. Since all routers are in contact with config servers, we need not add the shards again to router-2.*

## Inserting and Sharding Data

First let us create a database inside any one of the router with 2 collections in it.

```docker exec -it router-1 mongo```

```use bank``` - To create Database

```sh.enableSharding("bank")``` - To enable sharding for our database

```db.createCollection("users")``` - To create users collection

```db.createColeltion("transactions")``` - To create transactions collection

```sh.setBalancerState(true)``` - To set balancer state to true which ensures data balancing

```sh.startBalancer()``` - To run balancer at current time

```sh.shardCollection("bank.users",{"user_id":"hashed"});``` - To shard users collection

```sh.shardCollection("bank.transactions",{"transaction_id":"hashed"});``` - To shard transactions colletion

