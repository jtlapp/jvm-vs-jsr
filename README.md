# jvm-vs-jsr

Benchmarks comparing JVM and JS Runtime concurrency frameworks.

**CURRENTLY UNDER DEVELOPMENT**

## Introduction

When Node.js first came out, people were astonished that it could provide better throughput than 
Java for I/O-bound operations. Numerous reactive frameworks have since been created to rectify the 
issue in Java, but I had trouble finding apples-to-apples comparisons of how well these frameworks
perform relative to node.js. TechEmpower [provides benchmarks](https://www.techempower.com/benchmarks),
but I found it difficult to understand the commonalities and differences between any two
implementations. So I decided to create the present repo to definitively answer the question of how
Java presently compares to Node.js for I/O-bound work, though also as an exercise in helping me 
personally gain understanding. To provide a fuller comparison of modern Java to modern JS, I'll
also include benchmarks for Deno and Bun.

Please do let me know how I can improve the performance of any of these implementations.

## Installation and Setup

Install kubectl, helm, and helmfile, and configure kubectl for your cluster.

Run `helmfile init` to further install the Helm "diff" and "secrets" plugins.

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

## Building and Deploying

Create your cluster and configure `kubectl` to use it. Then:

```bash
mvn clean install
./bin/deploy database
./bin/deploy client
./bin/deploy spring-jdbc-kernel # or another app
```

The `deploy` command deploys or redeploys. In the case of apps, it replaces the currently 
deployed app (if any) with the named app.

## Running Benchmarks

1. Exec into the client pod using bash: `kubectl exec -it <client-pod> -- bash`.
2. Run `./benchmark setup <scenario>` to set up the scenario of the given name.
3. Run `./benchmark run <scenario> -rate <requests-per-sec> -duration <seconds>`.

Run `./benchmark` to see other commands and get usage help.

When running a test scenario, it outputs the first response for each unique combination
of status code, shared query name, and error message. For queries erroneously returning
non-JSON, it also prints each unique combination of status code and response body. This
output assists with debugging newly added applications.

## Useful Commands

```bash
./bin/undeploy database
./bin/undeploy client
./bin/undeploy spring-jdbc-kernel # or another app
```
