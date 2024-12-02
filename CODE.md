Follow the steps mentioned here to run the setup locally.

Before starting the application, make sure you have read the *SETUP.md* file to run code without any discrepancies.

## 1) Starting the Frontend Server
Navigate to UserInterface folder. Then run the following command:
```
http-server .
```
This will start a server that will the frontend files on a port. The port will be available on the terminal. Visit the URL in the terminal for UI of the application.

## 2) Starting MongoDB cluster server
Navigate to scripts folder. Then run **main.sh** script. Make sure the script has permissions to be run.

This script will start all the docker containers with 3 config-servers, 3 routers and 3 shard-replica sets (each set with 3 nodes).

## 3) Inserting Data
Visit the URL: https://drive.google.com/drive/folders/1BhCH97EacCjACIv4A69JWZjUbs4k6XWE
Download all the json files which contain information about 200,000 users and 1.5 million transactions.

To insert data, open MongoDB Compass and connect to any one of the router. For example, if the router is running on the port 27151, the url would be ```mongodb://localhost:27151/```. You can pick any of the 3 routers that were created.

After connecting to the router, click on **bank** database. It will have 2 collections namely *users* and *transactions*. 

To insert data, click on *Add Data* and select *Import JSON file*. Use corresponding json data files to insert into the collections. Please use *update_userInfo.json* for *user* collection and *mock_transactions.json* for transactions collection.

Since we are inserting data via GUI, follow the steps below to create indices on data, so data can be indexed and retrieved swiftly.

For transactions collection, create *asc* indices on *sender_id* and *receiver_id*.

For users collection, create index on *user_id* key.

## 4) Starting the Backend server

To start a server, run the *server.exe* executable with *-port* flag. To handle more load, run multiple server instances on different ports.

Eg: ```./server.exe -port 8080```

Now, navigate to frontend and explore the functionalities !!!
