# java-benchmarks

Benchmarks for backend Java frameworks

## Preparation

Install kubectl and helm, and configure kubectl for your cluster.

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

Label the three nodes as follows:

```bash
kubectl label nodes <node-1-name> kubernetes.io/hostname=client --overwrite
kubectl label nodes <node-2-name> kubernetes.io/hostname=app --overwrite
kubectl label nodes <node-3-name> kubernetes.io/hostname=database --overwrite
```

Add the Helm repos for Prometheus and Grafana:

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
```

## Installation

Create your cluster and configure `kubectl` to use it. Then:

```bash
mvn clean install
./bin/deploy common
./bin/deploy spring-jdbc-kernel
```

## Testing

1. Exec into the client pod using bash: `kubectl exec -it <client-pod> -- bash`.
2. Modify `/scripts/setup.ts` to set up the desired queries.
3. Run `deno -A setup.ts`.
4. Run the appropriate Lua benchmarking script. E.g.:

```bash
% wrk -t1 -c1 -d1s -s order-items/query.lua http://api-service:8080
```

Useful termination commands:

```bash
./bin/undeploy common
./bin/undeploy spring-jdbc-kernel
```
