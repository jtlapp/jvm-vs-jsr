# java-benchmarks

Benchmarks for backend Java frameworks

```bash

mvn clean install
kind create cluster
./bin/deploy common
./bin/deploy spring-kernel-threads
```

Testing:

```bash
> kubectl exec -it <pod> -- ash 
% curl -X GET localhost:8080/api/setup
Completed setup.
% curl -X GET "localhost:8080/api/select?user=1&order=1"
{...JSON...}
% curl -X GET "localhost:8080/api/update?user=1&order=1"
Updated.
% vi test.lua
% wrk -t1 -c1 -d1s -s test.lua http://backend-api-service:8080
```

Useful termination commands:

```bash
./bin/undeploy spring-kernel-threads
./bin/undeploy common
kind delete cluster
```
