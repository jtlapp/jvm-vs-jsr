# java-benchmarks

Benchmarks for backend Java frameworks

```bash

kind create cluster

cd spring-kernel-threads
docker build -t jtlapp/spring-kernel-threads .
kind load docker-image jtlapp/spring-kernel-threads:latest

kubectl apply -f deployment/postgres-deployment.yaml
kubectl apply -f deployment/app-deployment.yaml
kubectl port-forward service/backend-api-service 8080:8080
```

```bash
> curl -X GET localhost:8080/api/setup
Completed setup.
> curl -X GET "localhost:8080/api/read?user=1&order=1"

...

curl -X POST localhost:8080/actuator/shutdown
kind delete cluster
```
