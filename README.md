# jvm-vs-js

Benchmarks comparing JVM and JS Runtime concurrency frameworks.

**CURRENTLY UNDER DEVELOPMENT**

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
kubectl label nodes <node-3-name> kubernetes.io/hostname=backend --overwrite
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
./bin/jc-deploy database
./bin/jc-deploy client
./bin/jc-deploy spring-jdbc-kernel # or another app
```

## Testing

1. Exec into the client pod using bash: `kubectl exec -it <client-pod> -- bash`.
2. Run `./benchmark <test-suite> setup-all` to set up the test suite of the given name.
3. Run `./benchmark <test-suite> test -rate <requests-per-sec> -duration <seconds>`.

Run `./benchmark` to get usage help.

When running a test, the test outputs the first response for each unique combination of
status code, shared query name, and error message. For queries erroneously returning
non-JSON, it also prints each unique combination of status code and response body.

## Useful commands:

```bash
./bin/jc-redeploy <release-name>
./bin/jc-replace <deployed-release> <replacement-release>

./bin/jc-undeploy database
./bin/jc-undeploy client
./bin/jc-undeploy spring-jdbc-kernel # or another app
```

## Alternative Helm Deployments

I'm temporarily also using this repo to explore various approaches to deploying Helm charts to 
Kubernetes. The following approaches are supported:

### Just Charts

This approach deploys using only helm charts via the `bin/` scripts prefixed `jc-`.

### Helmfile

This approach deploys using helmfile and helm charts via the `bin/` scripts prefixed `hf-`. The 
approach works but is not yet idiomatic.

Run `hf-deploy` to deploy a chart or replace an existing chart in the same chart group. Run 
`hf-undeploy` to remove a chart.

Requires first installing `helmfile` and running `helmfile init` to further install the Helm "diff" 
and "secrets" plugins.

### Timoni

PLanned, but not yet supported. [See here.](https://timoni.sh/)
