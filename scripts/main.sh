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

echo -e "${CYAN}${BOLD}Starting MongoDB Sharded Cluster Setup...${RESET}"

# Step 1: Create Docker Network
echo -e "\n${GREEN}${BOLD}Step 1: Create Docker Network${RESET}"
echo -e "${BLUE}Creating Docker network...${RESET}"
./create_network.sh
echo -e "${GREEN}[Done]${RESET}"
echo -e "\n"  # Extra blank line for spacing

# Step 2: Setup Config Servers
echo -e "${GREEN}${BOLD}Step 2: Setup Config Servers${RESET}"
echo -e "${BLUE}Setting up Config Servers...${RESET}"
./config_servers.sh
echo -e "${GREEN}[Done]${RESET}"
echo -e "\n"  # Extra blank line for spacing

# Step 3: Setup Shard Replica Sets
echo -e "${GREEN}${BOLD}Step 3: Setup Shard Replica Sets${RESET}"
echo -e "${BLUE}Setting up Shard Replica Sets...${RESET}"
./create_shards.sh
echo -e "${GREEN}[Done]${RESET}"
echo -e "\n"  # Extra blank line for spacing

# Step 4: Setup Routers
echo -e "${GREEN}${BOLD}Step 4: Setup Routers${RESET}"
echo -e "${BLUE}Setting up Routers...${RESET}"
./create_routers.sh
echo -e "${GREEN}[Done]${RESET}"
echo -e "\n"  # Extra blank line for spacing

# Step 5: Connect Shards to Routers
echo -e "${GREEN}${BOLD}Step 5: Connect Shards to Routers${RESET}"
echo -e "${BLUE}Connecting Shards to Routers...${RESET}"
./connect_shards.sh
echo -e "${GREEN}[Done]${RESET}"
echo -e "\n"  # Extra blank line for spacing

echo -e "${CYAN}${BOLD}MongoDB Sharded Cluster Setup Completed Successfully!${RESET}"
echo -e "\n"  # Extra blank line for final touch
