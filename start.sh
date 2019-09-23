echo "Removing any old services:"
docker-compose rm -sf bootstrap-node kademlia-node
echo "Building app:"
go build -o ./bin/kademlia ./src/kademlia/
echo "Deploying bootstrap node:"
docker-compose up -d bootstrap-node
echo "Deploying kademlia nodes:"
docker-compose up -d --scale kademlia-node=10 kademlia-node