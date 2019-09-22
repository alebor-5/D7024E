echo "Removing any old services:"
docker-compose rm -sf bootstrap-node kademlia-node
echo "Building app:"
go build -o ./bin/kademlia ./src/kademlia/
echo "Deploying bootstrap node:"
docker-compose up -d bootstrap-node
echo "Deploying kademlia nodes:"
docker-compose up -d --scale kademlia-node=10 kademlia-node

# interval=1 #Number of seconds to wait
# n=10 #Number of tries

# docker stack rm kademlia
# echo "Building app"
# go build -o ./bin/kademlia ./src/kademlia/
# echo "Restarting Stack"
# for i in $(seq 0 $n); do 
#     (docker stack deploy -c docker-compose.yml kademlia) && break
#     if [ $i != $n ]; then
#         echo "Trying again in $interval second(s)...";
#         sleep $interval;
#     fi;
# done