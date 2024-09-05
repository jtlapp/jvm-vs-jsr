# java-benchmarks

Benchmarks for backend Java frameworks

```bash

mvn clean install
kind create cluster
./bin/deploy common
./bin/deploy spring-kernel-threads
kubectl port-forward service/backend-api-service 8080:8080
```

Testing:

```bash
> curl -X GET localhost:8080/api/setup
Completed setup.
> curl -X GET "localhost:8080/api/select?user=1&order=1"
{...JSON...}
> curl -X GET "localhost:8080/api/update?user=1&order=1"
Updated.
```

Useful termination commands:

```bash
./bin/undeploy spring-kernel-threads
pkill -f "kubectl port-forward"
curl -X POST localhost:8080/actuator/shutdown
kind delete cluster
```
