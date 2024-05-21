# minikube

## ingress  
https://kubernetes.io/docs/tasks/access-application-cluster/ingress-minikube/  
minikube addons enable ingress  
kubectl get pods -n ingress-nginx  
kubectl apply -f project-ingress.yaml
kubectl get ingress  
minikube ip  
minikube tunnel to run ingress on macos

## host file  
https://kinsta.com/knowledgebase/edit-mac-hosts-file/  
sudo nano /etc/hosts  
add 127.0.0.1 project.example  
sudo killall -HUP mDNSResponder  