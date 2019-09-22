# D7024E
D7024E Lab Assignment - P2P Distributed Data Store

# Setup
## Install dependencies
- docker
- docker-compose
- golang

## Build docker image
Run `./docker_setup.sh` in the project root.

# Deployment
1. Navigate to project root
2. Run `./start.sh`

## Aditional notes
Running `./start.sh` will **recompile** the source code. To redeploy the container services simply run `./start.sh` again. To stop and remove the container services run `docker-compose rm -sf bootstrap-node kademlia-node`