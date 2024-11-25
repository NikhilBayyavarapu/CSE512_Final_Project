#!/bin/bash

# Containers list
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

# Function to start containers
start_container() {
  container_name=$1
  echo -e "\033[1;32mStarting container: \033[1;34m$container_name\033[0m"
  # Replace this with your actual container start command, e.g., docker start $container_name
  docker start "$container_name"
  if [ $? -eq 0 ]; then
    echo -e "\033[1;32mSuccessfully started $container_name\033[0m"
  else
    echo -e "\033[1;31mFailed to start $container_name\033[0m"
  fi
}

# Start each container
echo -e "\033[1;36mStarting containers...\033[0m"
for container in "${containers[@]}"; do
  start_container "$container"
done

echo -e "\033[1;36mAll containers have been started.\033[0m"
