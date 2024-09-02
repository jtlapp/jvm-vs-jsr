# java-benchmarks

Benchmarks for backend Java frameworks

```bash

kind create cluster
./bin/deploy common

cd spring-kernel-threads
docker build -t jtlapp/spring-kernel-threads .
kind load docker-image jtlapp/spring-kernel-threads:latest

cd ..
./bin/deploy spring-kernel-threads
kubectl port-forward service/backend-api-service 8080:8080
```

```bash
> curl -X GET localhost:8080/api/setup
Completed setup.
> curl -X GET "localhost:8080/api/select?user=1&order=1"
{...JSON...}
> curl -X GET "localhost:8080/api/update?user=1&order=1"
Updated.
```

```
./bin/undeploy spring-kernel-threads
./bin/undeploy common

curl -X POST localhost:8080/actuator/shutdown
kind delete cluster
```
