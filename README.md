# java-benchmarks

Benchmarks for backend Java frameworks

## Preparation

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

Note: I found cross-platform building using `buildx` too unreliable to use, 
perhaps because of the dependency on calling `docker buildx create --use`. 
The client image therefore builds only for `amd64`.

## Installation

Create your cluster and configure `kubectl` to use it. Then:

```bash
mvn clean install
./bin/deploy common
./bin/deploy spring-kernel-threads
```

## Testing

1. Exec into the client pod using bash: `kubectl exec -it <clien-pod> -- bash`.
2. Modify `/scripts/setup.ts` to set up the desired queries.
3. Run the appropriate Lua benchmarking script. E.g.:

```bash
% wrk -t1 -c1 -d1s -s order-items/query.lua http://api-service:8080
```

Useful termination commands:

```bash
helm uninstall spring-kernel-threads
helm uninstall common
kind delete cluster
```
