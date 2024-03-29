nodes=${1:-50} # number of total nodes (including the bootstrap node) that should be deployed
echo "Removing any old services:"
docker-compose rm -sf bootstrap-node kademlia-node # -s, stops any running kademlia nodes. -f force removal of the containers.
echo "Building app:"
go build -o ./bin/kademlia ./src/kademlia/
echo "Deploying bootstrap node:"
docker-compose up -d bootstrap-node # -d, starts in detached mode
if [ $nodes -gt 1 ]
then
    echo "Deploying kademlia nodes:"
fi
i=$([ $(((nodes-1)%10)) == 0 ] && (echo 10) || (echo $(((nodes-1)%10))))
while [ $i -lt $(( nodes )) ]
do 
    docker-compose up -d --scale kademlia-node=$i kademlia-node # -d, starts in detached mode
    i=$(( i+10 ))
done
