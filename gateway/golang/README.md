# GATEWAY  

## compile grpc  
protoc user.proto --go_out=. --go-grpc_out=.  
protoc requestresponse.proto --go_out=. --go-grpc_out=.  
export PATH="$PATH:$(go env GOPATH)/bin"  

## build image  
docker build -t project-gateway:1.0.0 .  

## minikube  
kubectl get deployment project-gateway-golang  
kubectl expose deployment project-gateway-golang --type=NodePort --port=10001  
kubectl get service project-gateway-golang  
minikube service project-gateway-golang --url  