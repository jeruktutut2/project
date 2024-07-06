# USER  

## Run in docker  

## crate network  
docker network create project  
docker network ls  
if you want to use bridge network (default), you no need to create network  

## create mysql with network  
docker run --name project-mysql-network --network project -e MYSQL_ROOT_PASSWORD=12345 -p 3310:3306 -d mysql:8.1.0  
if you want to use bridge network (default), you no need to put --network project  

## create redis with network
docker run --name project-redis-network --network project -p 6380:6379 -d redis:latest
if you want to use bridge network (default), you no need to put --network project


## change host mysql  
MYSQL_HOST=project-mysql-network  
DATABASE_URL="project-mysql-network://username:password@host:3306/databasename"  
REDIS_HOST=project-redis-network  

## build image  
docker build -t project-backend-user:1.0.0 .  

## run project  
docker run --name project-user-container --network project -p 10001:10001 project-user:1.0.0
docker run --name project-user-container -p 10001:10001 project-user:1.0.0  
if you want to use bridge network (default), you no need to put --network project  

## test grpc golang
https://github.com/grpc/grpc-go/issues/1786  