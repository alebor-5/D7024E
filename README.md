# D7024E
## D7024E Lab Assignment - P2P Distributed Data Store

# Setup
## Install dependencies
You need to install the following dependencies:
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Golang](https://golang.org/)

## Build docker image
Before you can start the containers, you first need to build the docker image. As the image need to be named correctly we provide a script for you.
To execute the script simply run `./docker_setup.sh` in the project root.


# Deployment
To start the containers we provided a script called `start.sh`. This script starts one bootstrap node and 49 “regular” nodes. All the regular nodes will send a FIND_NODE RPC to the boostrap node with their own KademliaID. The bootstrap node will have a fixed KademliaID and also a fixed IP address (*10.0.0.2*).
To execute the startup script run `./start.sh` in the project root. 
If you whish to change the amount of nodes you can add an argument to script, for example `./start.sh 100` will start 100 nodes. 


## Aditional notes
Running `./start.sh` will **recompile** the source code. To redeploy the container services simply run `./start.sh` again. To stop and remove the container services run `docker-compose rm -sf bootstrap-node kademlia-node`

# CLI

There are currently 6 commands that can be used to interact with this Kademlia implementation, these are:

`Help` – Provides information about available command.

`Put [string]` – Uploads [string] to the Kademlia network, outputs the hash of the uploaded string.

`Get [hash]` – Outputs object corresponding to the [hash] if the object exists in the network.

`Exit` – Terminates the node.

`Enablelog` – Enables logging

`Disablelog` – Disables logging
