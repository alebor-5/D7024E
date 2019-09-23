nodes=10 # number of total nodes (including the bootstrap node) that should be deployed
echo "Removing any old services:"
docker-compose rm -sf bootstrap-node kademlia-node # -s, stops any running kademlia nodes. -f force removal of the containers.
echo "Building app:"
go build -o ./bin/kademlia ./src/kademlia/
echo "Deploying bootstrap node:"
docker-compose up -d bootstrap-node # -d, starts in detached mode
echo "Deploying kademlia nodes:"
docker-compose up -d --scale kademlia-node=($nodes-1) kademlia-node # -d, starts in detached mode