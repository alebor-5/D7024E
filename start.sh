interval=1 #Number of seconds to wait
n=10 #Number of tries

docker stack rm kademlia
echo "Building app"
go build -o ./bin/kademlia ./src/kademlia/
echo "Restarting Stack"
for i in $(seq 0 $n); do 
    (docker stack deploy -c docker-compose.yml kademlia) && break
    if [ $i < $n ]; then
        echo "Trying again in $interval second(s)...";
        sleep $interval;
    fi;
done