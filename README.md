# jvm-vs-jsr

Benchmarks comparing JVM and JS Runtime concurrency frameworks.

**CURRENTLY UNDER DEVELOPMENT**

## Purpose

To acquire a definitive understanding of how the throughput of I/O-bound work on Java virtual
threads compares with Java kernel threads, Java reactive frameworks, and JavaScript runtime
platforms (e.g. Node.js). Only PostgreSQL queries are examined.

## Introduction

When Node.js first came out, people were astonished that it could provide better throughput than
Java for I/O-bound operations. Numerous reactive frameworks have since been created to rectify the
issue (Java), but I had trouble finding apples-to-apples comparisons of how well these frameworks
perform relative to Node.js. TechEmpower [provides benchmarks](https://www.techempower.
com/benchmarks),
but I found it difficult to understand the commonalities and differences between any two
implementations. The present repo is an effort to get a definitive sense of how throughput compares
on these platforms for I/O-bound work, specifically PostgreSQL queries.

I'll be collecting throughput metrics for non-blocking APIs and for APIs that hit a backend
Postgres database. I'll be considering the
database clients JDBC (including HikariCP), R2DBC, and the Vert.x PG client.

I ended up writing an unexpected amount of Go client code for this project, but there are no
unit tests because I'm not planning to maintain the repo long term.

## The Plan

The first goal is to compare bare-metal Java solutions to bare-metal JavaScript runtime
solutions to learn the best possible throughput available. Fastify seems like a
reasonable bare-bones approach to testing on a JavaScript runtime, while Jooby seems like a
reasonable bare-bones approach for the JVM, given that Jooby code can be written in an
Express.js-like fashion. I also want to compare these results to Spring Boot, which seems like a
necessary baseline for comparing the performance of Java frameworks. I'll then look at how
Quarkus and Micronaut perform with R2DBC reactive I/O.

Here are the JavaScript Runtime combinations I'll be benchmarking:

- Fastify
  - Clustered Node.js + postgresjs
  - Clustered Deno + postgresjs
  - Clustered Bun + postgresjs

Here are the JVM combinations I'll be benchmarking:

- Jooby
  - JDBC with kernel threads (Java)
  - JDBC with virtual threads (Java)
  - R2DBC (Java)
  - Vert.x (Java)
- Spring Boot
  - JDBC with kernel threads (Java)
  - JDBC with virtual threads (Java)
  - WebFlux/R2DBC (Java)
  - WebFlux/R2DBC (Kotlin)
- Quarkus
  - R2DBC (Kotlin)
- Micronaut
  - R2DBC (Kotlin)

I'm benchmarking the reactive framework combinations in Kotlin because I find it unreasonable to
implement reactive frameworks for the JVM without coroutines, for the same reason it's
unreasonable to implement I/O in Node.js without async/await.

For each combination, I'll also be implemented a thread sleep to measure raw concurrency
unimpeded by waiting on a third tier.

## Installation and Setup

Install kubectl, helm, and helmfile, and configure kubectl for your cluster.

Run `helmfile init` to further install the Helm "diff" and "secrets" plugins.

Set the docker hostname image prefix in the following environment variable, excluding any trailing
+`/`:

```bash
export DOCKER_IMAGE_PREFIX=<your-docker-image-prefix>
```

And set it in `charts/values-local.yaml` (which is in `.gitignore`):

```bash
global:
  app:
    dockerImagePrefix: <your-docker-image-prefix>
```

If your cluster does not automatically provision persistent volumes, you'll also need to create 
a volume called `client-postgres-volume`.

Note: I found cross-platform building using `buildx` too unreliable to use,
perhaps because of the dependency on calling `docker buildx create --use`.
The client image therefore builds only for `amd64`.

Label the three nodes as follows. You may need to run `kubectl get nodes` to list the node names.

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
./bin/deploy backend-database
./bin/deploy client
./bin/deploy spring-jdbc-kernel-app # or another app
```

The `deploy` command deploys or redeploys. In the case of apps, it replaces the currently
deployed app (if any) with the named app.

## Running Benchmarks

1. Exec into the client pod using bash: `kubectl exec -it <client-pod> -- bash`.
2. Run `./bencjmark setup-results` to create the benchmark results database.
3. Run `./benchmark setup-backend -scenario <scenario>` to set up the scenario of the given name.
   Only required if the scenario uses backend Postgres tables.
4. 
5. TODO: introduce try, run, and loop
6. Run `./benchmark run -scenario <scenario> -rate <requests-per-sec> -duration <seconds>`.

Run `./benchmark` to see other commands and get usage help.

When running a test scenario, it outputs the first response for each unique combination
of status code, shared query name, and error message. For queries erroneously returning
non-JSON, it also prints each unique combination of status code and response body. This
output assists with debugging newly added applications.

## Useful Commands

```bash
./bin/undeploy backend-database
./bin/undeploy client
./bin/undeploy spring-jdbc-kernel-app # or another app
```

## Using Minikube (Couldn't get to work)

Do the following to set up benchmarks for Minikube. Mind you, Minikube doesn't let you specify 
resource allocations per node.

- If you've previously run minikube and have an existing cluster, you won't be able to set the 
  number of CPUs that the cluster should use unless you first run `minikube delete`.
- Start minikube indicating the number of CPUs you'd like it to use across all nodes. E.g. 
  `minikube start --cpus 6 --insecure-registry "10.0.0.0/24"`. This restricts the docker VM to 
  using this many CPU units.
- Now have minikube create a docker image registry: `minikube addons enable registry`. Ignore 
  the port that minikube reports.
- In a dedicated shell, run `docker run --rm -it --network=host alpine ash -c "apk add socat && socat TCP-LISTEN:5000,reuseaddr,fork TCP:$(minikube ip):5000"
` to establish the registry within the VM. This command blocks the terminal.
- Use `localhost:5000` as your docker image prefix, according to the above instructions.
- Run `mvn clean install` to build the images with this image prefix. (If you want to first 
  remove all prior images, run `docker image prune -a`.)
- Create three nodes by running `minikube node add` three times (e.g. `minikube node add && 
  minikube node add && minikube node add`).
- Label the three nodes other than the control-plane node per the above instructions.

## Using Kind

Do the following to perform benchmarks on a Kind cluster. Mind you, Kind doesn't let you specify
resource allocations per node.

- Create a Kind cluster with three worker nodes via `kind create cluster --config 
  config/kind-config.yaml`. Note that this will create a directory called `client-pv` in the 
  current  directory to hold the client's persistent volume. `client-pv/` is in `.gitignore`. 
  This also labels the nodes for you, so you don't have to label them manually.
- Following the above instructions, set the image prefix to any non-empty string.

## TODO

- Remove image prefix
- Remove value-local.yaml
- Hard-code scripts for Kind
- 