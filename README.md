# D7024E
## D7024E Lab Assignment - P2P Distributed Data Store

# Setup
## Install dependencies
You need to install the following dependencies:
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Golang](https://golang.org/)

## Build docker image
Before you can start the containers, you first need to build the docker image. As the image need to be named correctly we provided a script for you.
To execute the script simply run `./docker_setup.sh` in the project root.


# Deployment
To start the containers we provided a script called `start.sh`. This script starts one bootstrap node and 49 “regular” nodes. All the regular nodes will run the join procedure with the bootstrap nodes static ID and IP (*10.0.0.2*).
To execute the startup script run `./start.sh` in the project root. 
If you whish to change the amount of nodes you can add an argument to script, for example `./start.sh 20` will start 20 nodes. 


## Aditional notes
Running `./start.sh` will **recompile** the source code. To redeploy the container services simply run `./start.sh` again. To stop and remove the container services run `docker-compose rm -sf bootstrap-node kademlia-node`

# Test
To run all tests in the project. Navigate to D7024E/src/kademlia/ then run `go test -cover`.

# CLI

There are currently 6 commands that can be used to interact with this Kademlia implementation, these are:

`Help` – Provides information about available command.

`Put [string]` – Uploads [string] to the Kademlia network, outputs the hash of the uploaded string.

`Get [hash]` – Outputs object corresponding to the [hash] if the object exists in the network.

`Exit` – Terminates the node.

`Enablelog` – Enables logging

`Disablelog` – Disables logging

# RESTapi
All nodes can handle http requests from the internal kademlia network on port *8080*. The bootstrap nodes port *8080* is mapped to the host machines port *10000*. Requests sent to *localhost:10000* will be forwarded to the bootstrap node.

### Supported requests
- **POST /objects**
    - Stores the data provided in the request body
    - A succesfull call responds with *201 CREATED*
    - The response contains a header specifying the location of the stored data *Location: /objects/{hash}*
    - The response body is the same as the request body
- **GET /objects/{hash}**
    - If a valid hash key is provided it returns the corresponding value in the response body

*{hash} is a hexadecimal value consisting of 40 characters.*
