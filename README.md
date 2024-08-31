# java-benchmarks

Benchmarks for backend Java frameworks

```bash

kind create cluster

cd spring-kernel-threads
docker build -t jtlapp/spring-kernel-threads .
kind load docker-image jtlapp/spring-kernel-threads:latest

kubectl apply -f deployment/init-db-configmap.yaml
kubectl apply -f deployment/postgres-deployment.yaml
kubectl apply -f deployment/app-deployment.yaml
kubectl port-forward service/backend-api-service 8080:8080
```

```bash
> curl -X POST -d "First" localhost:8080/api/message
1
> curl -X GET -d "First" localhost:8080/api/message/1
First

...

curl -X POST localhost:8080/actuator/shutdown
kind delete cluster
```
