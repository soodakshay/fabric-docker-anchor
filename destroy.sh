docker kill $(docker ps -q)
docker rm $(docker ps -qa)
echo y | docker volume prune
echo y | docker network prune
