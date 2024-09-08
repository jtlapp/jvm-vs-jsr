# java-benchmarks

Benchmarks for backend Java frameworks

## Preparation

The root POM uses `docker buildx` to build both amd64 and arm64 images. 
You'll need to set this up as follows:

```bash
docker buildx create --use
```

Set the docker image prefix in an environment variables:

```bash
export DOCKER_IMAGE_PREFIX=<your-docker-image-prefix>
```

And set it in `charts/values-secret.yaml` (which is in `.gitignore`):

```bash
global:
  app:
    dockerImagePrefix: <your-docker-image-prefix>
```

## Installation

Create your cluster and configure `kubectl` to use it. Then:

```bash
mvn clean install
./bin/deploy common
./bin/deploy spring-kernel-threads
```

Testing:

```bash
> kubectl exec -it <pod> -- ash 
% curl -X GET api-service:8080/api/setup
Completed setup.
% curl -X GET "api-service:8080/api/select?user=1&order=1"
{...JSON...}
% curl -X GET "api-service:8080/api/update?user=1&order=1"
Updated.
% vi test.lua
% wrk -t1 -c1 -d1s -s test.lua http://api-service:8080
```

Useful termination commands:

```bash
./bin/undeploy spring-kernel-threads
./bin/undeploy common
kind delete cluster
```
